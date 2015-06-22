package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bitly/go-nsq"
	"github.com/dancannon/gorethink"
	"github.com/johntdyer/slackrus"
	//"github.com/pzduniak/glogrus"
	"github.com/getsentry/raven-go"
	"github.com/willf/bloom"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"gopkg.in/igm/sockjs-go.v2/sockjs"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/db"
	"github.com/lavab/api/env"
	"github.com/lavab/api/routes"
	"github.com/lavab/api/utils"
)

// sessions contains all "subscribing" WebSockets sessions
var (
	sessions     = map[string][]sockjs.Session{}
	sessionsLock sync.Mutex
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// PrepareMux sets up the API
func PrepareMux(flags *env.Flags) *web.Mux {
	// Set up a new logger
	log := logrus.New()

	// Set the formatter depending on the passed flag's value
	if flags.LogFormatterType == "text" {
		log.Formatter = &logrus.TextFormatter{
			ForceColors: flags.ForceColors,
		}
	} else if flags.LogFormatterType == "json" {
		log.Formatter = &logrus.JSONFormatter{}
	}

	// Install Logrus hooks
	if flags.SlackURL != "" {
		var level []logrus.Level

		switch flags.SlackLevels {
		case "debug":
			level = slackrus.LevelThreshold(logrus.DebugLevel)
		case "error":
			level = slackrus.LevelThreshold(logrus.ErrorLevel)
		case "fatal":
			level = slackrus.LevelThreshold(logrus.FatalLevel)
		case "info":
			level = slackrus.LevelThreshold(logrus.InfoLevel)
		case "panic":
			level = slackrus.LevelThreshold(logrus.PanicLevel)
		case "warn":
			level = slackrus.LevelThreshold(logrus.WarnLevel)
		}

		log.Hooks.Add(&slackrus.SlackrusHook{
			HookURL:        flags.SlackURL,
			AcceptedLevels: level,
			Channel:        flags.SlackChannel,
			IconEmoji:      flags.SlackIcon,
			Username:       flags.SlackUsername,
		})
	}

	// Connect to raven
	var rc *raven.Client
	if flags.RavenDSN != "" {
		h, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}

		rc, err = raven.NewClient(flags.RavenDSN, map[string]string{
			"hostname": h,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	env.Raven = rc

	// Pass it to the environment package
	env.Log = log

	// Load the bloom filter
	bf := bloom.NewWithEstimates(flags.BloomCount, 0.001)
	bff, err := os.Open(flags.BloomFilter)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unable to open the bloom filter file")
	}
	defer bff.Close()
	if _, err := bf.ReadFrom(bff); err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unable to read from the bloom filter file")
	}
	env.PasswordBF = bf

	// Initialize the cache
	redis, err := cache.NewRedisCache(&cache.RedisCacheOpts{
		Address:  flags.RedisAddress,
		Database: flags.RedisDatabase,
		Password: flags.RedisPassword,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to connect to the redis server")
	}

	env.Cache = redis

	// Set up the database
	rethinkOpts := gorethink.ConnectOpts{
		Address:  flags.RethinkDBAddress,
		AuthKey:  flags.RethinkDBKey,
		Database: flags.RethinkDBDatabase,
		MaxIdle:  10,
		Timeout:  time.Second * 10,
	}
	err = db.Setup(rethinkOpts)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to set up the database")
	}

	// Initialize the actual connection
	rethinkOpts.Database = flags.RethinkDBDatabase
	rethinkSession, err := gorethink.Connect(rethinkOpts)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to connect to the database")
	}

	// Put the RethinkDB session into the environment package
	env.Rethink = rethinkSession

	// Initialize the tables
	env.Tokens = &db.TokensTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"tokens",
		),
		Cache: redis,
	}
	env.Accounts = &db.AccountsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"accounts",
		),
		Tokens: env.Tokens,
	}
	env.Addresses = &db.AddressesTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"addresses",
		),
	}
	env.Keys = &db.KeysTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"keys",
		),
	}
	env.Contacts = &db.ContactsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"contacts",
		),
	}
	env.Reservations = &db.ReservationsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"reservations",
		),
	}
	env.Emails = &db.EmailsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"emails",
		),
	}
	env.Threads = &db.ThreadsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"threads",
		),
	}
	env.Labels = &db.LabelsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"labels",
		),
		Emails: env.Emails,
		//Cache:  redis,
	}
	env.Files = &db.FilesTable{
		Emails: env.Emails,
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"files",
		),
	}

	// Create a producer
	producer, err := nsq.NewProducer(flags.NSQdAddress, nsq.NewConfig())
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unable to create a new nsq producer")
	}

	/*defer func(producer *nsq.Producer) {
		producer.Stop()
	}(producer)*/

	env.Producer = producer

	// Get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unable to get the hostname")
	}

	// Create a delivery consumer
	deliveryConsumer, err := nsq.NewConsumer("email_delivery", hostname, nsq.NewConfig())
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"topic": "email_delivery",
		}).Fatal("Unable to create a new nsq consumer")
	}
	//defer deliveryConsumer.Stop()

	deliveryConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
		// Raven recoverer
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			msg := &raven.Message{
				Message: string(m.Body),
				Params:  []interface{}{"delivery"},
			}

			var packet *raven.Packet
			switch rval := recover().(type) {
			case error:
				packet = raven.NewPacket(rval.Error(), msg, raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
			default:
				str := fmt.Sprintf("%+v", rval)
				packet = raven.NewPacket(str, msg, raven.NewException(errors.New(str), raven.NewStacktrace(2, 3, nil)))
			}

			rc.Capture(packet, nil)
		}()

		var msg *struct {
			ID    string `json:"id"`
			Owner string `json:"owner"`
		}

		if err := json.Unmarshal(m.Body, &msg); err != nil {
			return err
		}

		// Check if we are handling owner's session
		if _, ok := sessions[msg.Owner]; !ok {
			return nil
		}

		if len(sessions[msg.Owner]) == 0 {
			return nil
		}

		// Resolve the email
		email, err := env.Emails.GetEmail(msg.ID)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    msg.ID,
			}).Error("Unable to resolve an email from queue")
			return nil
		}

		// Resolve the thread
		thread, err := env.Threads.GetThread(email.Thread)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":  err.Error(),
				"id":     msg.ID,
				"thread": email.Thread,
			}).Error("Unable to resolve a thread from queue")
			return nil
		}

		// Send notifications to subscribers
		for _, session := range sessions[msg.Owner] {
			result, _ := json.Marshal(map[string]interface{}{
				"type":   "delivery",
				"id":     msg.ID,
				"name":   email.Name,
				"thread": email.Thread,
				"labels": thread.Labels,
			})
			err = session.Send(string(result))
			if err != nil {
				env.Log.WithFields(logrus.Fields{
					"id":    session.ID(),
					"error": err.Error(),
				}).Warn("Error while writing to a WebSocket")
			}
		}

		return nil
	}), 10)

	if err := deliveryConsumer.ConnectToNSQLookupd(flags.LookupdAddress); err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unable to connect to nsqlookupd")
	}

	// Create a receipt consumer
	receiptConsumer, err := nsq.NewConsumer("email_receipt", hostname, nsq.NewConfig())
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"topic": "email_receipt",
		}).Fatal("Unable to create a new nsq consumer")
	}
	//defer receiptConsumer.Stop()

	receiptConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
		// Raven recoverer
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			msg := &raven.Message{
				Message: string(m.Body),
				Params:  []interface{}{"receipt"},
			}

			var packet *raven.Packet
			switch rval := recover().(type) {
			case error:
				packet = raven.NewPacket(rval.Error(), msg, raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
			default:
				str := fmt.Sprintf("%+v", rval)
				packet = raven.NewPacket(str, msg, raven.NewException(errors.New(str), raven.NewStacktrace(2, 3, nil)))
			}

			rc.Capture(packet, nil)
		}()

		var msg *struct {
			ID    string `json:"id"`
			Owner string `json:"owner"`
		}

		if err := json.Unmarshal(m.Body, &msg); err != nil {
			return err
		}

		// Check if we are handling owner's session
		if _, ok := sessions[msg.Owner]; !ok {
			return nil
		}

		if len(sessions[msg.Owner]) == 0 {
			return nil
		}

		// Resolve the email
		email, err := env.Emails.GetEmail(msg.ID)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    msg.ID,
			}).Error("Unable to resolve an email from queue")
			return nil
		}

		// Resolve the thread
		thread, err := env.Threads.GetThread(email.Thread)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":  err.Error(),
				"id":     msg.ID,
				"thread": email.Thread,
			}).Error("Unable to resolve a thread from queue")
			return nil
		}

		// Send notifications to subscribers
		for _, session := range sessions[msg.Owner] {
			result, _ := json.Marshal(map[string]interface{}{
				"type":   "receipt",
				"id":     msg.ID,
				"name":   email.Name,
				"thread": email.Thread,
				"labels": thread.Labels,
			})
			err = session.Send(string(result))
			if err != nil {
				env.Log.WithFields(logrus.Fields{
					"id":    session.ID(),
					"error": err.Error(),
				}).Warn("Error while writing to a WebSocket")
			}
		}

		return nil
	}), 10)

	if err := receiptConsumer.ConnectToNSQLookupd(flags.LookupdAddress); err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unable to connect to nsqlookupd")
	}

	// Create a new goji mux
	mux := web.New()

	// Include the most basic middlewares:
	//  - RequestID assigns an unique ID for each request in order to identify errors.
	//  - Glogrus logs each request
	//  - Recoverer prevents panics from crashing the API
	//  - AutomaticOptions automatically responds to OPTIONS requests
	mux.Use(func(c *web.C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// sockjs doesn't want to work with our code, as the author doesn't understand http.Headers
			if strings.HasPrefix(r.RequestURI, "/ws") {
				h.ServeHTTP(w, r)
				return
			}

			// because why not
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			allowedHeaders := []string{
				"Origin",
				"Content-Type",
				"Authorization",
				"X-Requested-With",
			}

			reqHeaders := strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",")
			allowedHeaders = append(allowedHeaders, reqHeaders...)

			resultHeaders := []string{}
			seenHeaders := map[string]struct{}{}
			for _, val := range allowedHeaders {
				if _, ok := seenHeaders[val]; !ok && val != "" {
					resultHeaders = append(resultHeaders, val)
					seenHeaders[val] = struct{}{}
				}
			}

			w.Header().Set("Access-Control-Allow-Headers", strings.Join(resultHeaders, ","))

			/*
				if c.Env != nil {
					if v, ok := c.Env[web.ValidMethodsKey]; ok {
						if methods, ok := v.([]string); ok {
							methodsString := strings.Join(methods, ",")
							w.Header().Set("Allow", methodsString)
							w.Header().Set("Access-Control-Allow-Methods", methodsString)
						}
					}
				} */

			// yolo
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			if r.Method != "OPTIONS" {
				h.ServeHTTP(w, r)
			}
		})
	})
	mux.Use(middleware.RequestID)
	//mux.Use(glogrus.NewGlogrus(log, "api"))
	mux.Use(recoverer)
	mux.Use(middleware.AutomaticOptions)

	// Set up an auth'd mux
	auth := web.New()
	auth.Use(routes.AuthMiddleware)

	// Index route
	mux.Get("/", routes.Hello)

	// Accounts
	auth.Get("/accounts", routes.AccountsList)
	mux.Post("/accounts", routes.AccountsCreate)
	auth.Get("/accounts/:id", routes.AccountsGet)
	auth.Put("/accounts/:id", routes.AccountsUpdate)
	auth.Delete("/accounts/:id", routes.AccountsDelete)
	auth.Post("/accounts/:id/wipe-data", routes.AccountsWipeData)
	auth.Post("/accounts/:id/start-onboarding", routes.AccountsStartOnboarding)

	// Addresses
	auth.Get("/addresses", routes.AddressesList)

	// Avatars
	mux.Get(regexp.MustCompile(`/avatars/(?P<hash>[\S\s]*?)\.(?P<ext>svg|png)(?:[\S\s]*?)$`), routes.Avatars)
	//mux.Get("/avatars/:hash.:ext", routes.Avatars)

	// Files
	auth.Get("/files", routes.FilesList)
	auth.Post("/files", routes.FilesCreate)
	auth.Get("/files/:id", routes.FilesGet)
	auth.Put("/files/:id", routes.FilesUpdate)
	auth.Delete("/files/:id", routes.FilesDelete)

	// Tokens
	auth.Get("/tokens", routes.TokensGet)
	auth.Get("/tokens/:id", routes.TokensGet)
	mux.Post("/tokens", routes.TokensCreate)
	auth.Delete("/tokens", routes.TokensDelete)
	auth.Delete("/tokens/:id", routes.TokensDelete)

	// Threads
	auth.Get("/threads", routes.ThreadsList)
	auth.Get("/threads/:id", routes.ThreadsGet)
	auth.Put("/threads/:id", routes.ThreadsUpdate)
	auth.Delete("/threads/:id", routes.ThreadsDelete)

	// Emails
	auth.Get("/emails", routes.EmailsList)
	auth.Post("/emails", routes.EmailsCreate)
	auth.Get("/emails/:id", routes.EmailsGet)
	auth.Delete("/emails/:id", routes.EmailsDelete)

	// Labels
	auth.Get("/labels", routes.LabelsList)
	auth.Post("/labels", routes.LabelsCreate)
	auth.Get("/labels/:id", routes.LabelsGet)
	auth.Put("/labels/:id", routes.LabelsUpdate)
	auth.Delete("/labels/:id", routes.LabelsDelete)

	// Contacts
	auth.Get("/contacts", routes.ContactsList)
	auth.Post("/contacts", routes.ContactsCreate)
	auth.Get("/contacts/:id", routes.ContactsGet)
	auth.Put("/contacts/:id", routes.ContactsUpdate)
	auth.Delete("/contacts/:id", routes.ContactsDelete)

	// Keys
	mux.Get("/keys", routes.KeysList)
	auth.Post("/keys", routes.KeysCreate)
	mux.Get("/keys/:id", routes.KeysGet)
	auth.Post("/keys/:id/vote", routes.KeysVote)

	// Headers proxy
	mux.Get("/headers", func(w http.ResponseWriter, r *http.Request) {
		utils.JSONResponse(w, 200, r.Header)
	})

	mux.Handle("/ws/*", sockjs.NewHandler("/ws", sockjs.DefaultOptions, func(session sockjs.Session) {
		var subscribed string

		// A new goroutine seems to be spawned for each new session
		for {
			// Read a message from the input
			msg, err := session.Recv()
			if err != nil {
				if err != sockjs.ErrSessionNotOpen {
					env.Log.WithFields(logrus.Fields{
						"id":    session.ID(),
						"error": err.Error(),
					}).Warn("Error while reading from a WebSocket")
				}
				break
			}

			// Decode the message
			var input struct {
				Type    string            `json:"type"`
				Token   string            `json:"token"`
				ID      string            `json:"id"`
				Method  string            `json:"method"`
				Path    string            `json:"path"`
				Body    string            `json:"body"`
				Headers map[string]string `json:"headers"`
			}
			err = json.Unmarshal([]byte(msg), &input)
			if err != nil {
				// Return an error response
				resp, _ := json.Marshal(map[string]interface{}{
					"type":  "error",
					"error": err,
				})
				err := session.Send(string(resp))
				if err != nil {
					env.Log.WithFields(logrus.Fields{
						"id":    session.ID(),
						"error": err.Error(),
					}).Warn("Error while writing to a WebSocket")
					break
				}
				continue
			}

			// Check message's type
			if input.Type == "subscribe" {
				// Listen to user's events

				// Check if token is empty
				if input.Token == "" {
					// Return an error response
					resp, _ := json.Marshal(map[string]interface{}{
						"type":  "error",
						"error": "Invalid token",
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						break
					}
					continue
				}

				// Check the token in database
				token, err := env.Tokens.GetToken(input.Token)
				if err != nil {
					// Return an error response
					resp, _ := json.Marshal(map[string]interface{}{
						"type":  "error",
						"error": "Invalid token",
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						break
					}
					continue
				}

				// Do the actual subscription
				subscribed = token.Owner
				sessionsLock.Lock()

				// Sessions map already contains this owner
				if _, ok := sessions[token.Owner]; ok {
					sessions[token.Owner] = append(sessions[token.Owner], session)
				} else {
					// We have to allocate a new slice
					sessions[token.Owner] = []sockjs.Session{session}
				}

				// Unlock the map write
				sessionsLock.Unlock()

				// Return a response
				resp, _ := json.Marshal(map[string]interface{}{
					"type": "subscribed",
				})
				err = session.Send(string(resp))
				if err != nil {
					env.Log.WithFields(logrus.Fields{
						"id":    session.ID(),
						"error": err.Error(),
					}).Warn("Error while writing to a WebSocket")
					break
				}
			} else if input.Type == "unsubscribe" {
				if subscribed == "" {
					resp, _ := json.Marshal(map[string]interface{}{
						"type":  "error",
						"error": "Not subscribed",
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						break
					}
				}

				sessionsLock.Lock()

				if _, ok := sessions[subscribed]; !ok {
					// Return a response
					resp, _ := json.Marshal(map[string]interface{}{
						"type": "unsubscribed",
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						sessionsLock.Unlock()
						subscribed = ""
						break
					}
					sessionsLock.Unlock()
					subscribed = ""
					continue
				}

				if len(sessions[subscribed]) == 1 {
					delete(sessions, subscribed)

					// Return a response
					resp, _ := json.Marshal(map[string]interface{}{
						"type": "unsubscribed",
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						sessionsLock.Unlock()
						subscribed = ""
						break
					}
					sessionsLock.Unlock()
					subscribed = ""
					continue
				}

				// Find the session
				index := -1
				for i, session2 := range sessions[subscribed] {
					if session == session2 {
						index = i
						break
					}
				}

				// We didn't find anything
				if index == -1 {
					// Return a response
					resp, _ := json.Marshal(map[string]interface{}{
						"type": "unsubscribed",
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						sessionsLock.Unlock()
						subscribed = ""
						break
					}
					sessionsLock.Unlock()
					subscribed = ""
					continue
				}

				// We found it, so we are supposed to slice it
				sessions[subscribed][index] = sessions[subscribed][len(sessions[subscribed])-1]
				sessions[subscribed][len(sessions[subscribed])-1] = nil
				sessions[subscribed] = sessions[subscribed][:len(sessions[subscribed])-1]

				// Return a response
				resp, _ := json.Marshal(map[string]interface{}{
					"type": "unsubscribed",
				})
				err := session.Send(string(resp))
				if err != nil {
					env.Log.WithFields(logrus.Fields{
						"id":    session.ID(),
						"error": err.Error(),
					}).Warn("Error while writing to a WebSocket")
					sessionsLock.Unlock()
					subscribed = ""
					break
				}
				sessionsLock.Unlock()
				subscribed = ""
			} else if input.Type == "request" {
				// Perform the request
				w := httptest.NewRecorder()
				r, err := http.NewRequest(strings.ToUpper(input.Method), "http://api.lavaboom.io"+input.Path, strings.NewReader(input.Body))
				if err != nil {
					env.Log.WithFields(logrus.Fields{
						"id":    session.ID(),
						"error": err.Error(),
						"path":  input.Path,
					}).Warn("SockJS request error")

					// Return an error response
					resp, _ := json.Marshal(map[string]interface{}{
						"error": err.Error(),
					})
					err := session.Send(string(resp))
					if err != nil {
						env.Log.WithFields(logrus.Fields{
							"id":    session.ID(),
							"error": err.Error(),
						}).Warn("Error while writing to a WebSocket")
						break
					}
					continue
				}

				r.Body = nopCloser{strings.NewReader(input.Body)}

				r.RequestURI = input.Path

				for key, value := range input.Headers {
					r.Header.Set(key, value)
				}

				mux.ServeHTTP(w, r)

				// Return the final response
				result, _ := json.Marshal(map[string]interface{}{
					"type":    "response",
					"id":      input.ID,
					"status":  w.Code,
					"headers": w.HeaderMap,
					"body":    w.Body.String(),
				})
				err = session.Send(string(result))
				if err != nil {
					env.Log.WithFields(logrus.Fields{
						"id":    session.ID(),
						"error": err.Error(),
					}).Warn("Error while writing to a WebSocket")
					break
				}
			}
		}

		// We have to clear the subscription here too. TODO: make the code shorter
		if subscribed == "" {
			return
		}

		sessionsLock.Lock()

		if _, ok := sessions[subscribed]; !ok {
			sessionsLock.Unlock()
			return
		}

		if len(sessions[subscribed]) == 1 {
			delete(sessions, subscribed)
			sessionsLock.Unlock()
			return
		}

		// Find the session
		index := -1
		for i, session2 := range sessions[subscribed] {
			if session == session2 {
				index = i
				break
			}
		}

		// We didn't find anything
		if index == -1 {
			sessionsLock.Unlock()
			return
		}

		// We found it, so we are supposed to slice it
		sessions[subscribed][index] = sessions[subscribed][len(sessions[subscribed])-1]
		sessions[subscribed][len(sessions[subscribed])-1] = nil
		sessions[subscribed] = sessions[subscribed][:len(sessions[subscribed])-1]

		// Unlock the mutex
		sessionsLock.Unlock()
	}))

	// Merge the muxes
	mux.Handle("/*", auth)

	// Compile the routes
	mux.Compile()

	return mux
}
