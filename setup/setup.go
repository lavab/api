package setup

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
	"github.com/dancannon/gorethink"
	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"gopkg.in/igm/sockjs-go.v2/sockjs"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/db"
	"github.com/lavab/api/env"
	"github.com/lavab/api/factor"
	"github.com/lavab/api/routes"
	"github.com/lavab/glogrus"
)

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

	// Pass it to the environment package
	env.Log = log

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
		Address:     flags.RethinkDBAddress,
		AuthKey:     flags.RethinkDBKey,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
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

	// NATS queue connection
	nc, err := nats.Connect(flags.NATSAddress)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error":   err,
			"address": flags.NATSAddress,
		}).Fatal("Unable to connect to NATS")
	}

	c, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error":   err,
			"address": flags.NATSAddress,
		}).Fatal("Unable to initialize a JSON NATS connection")
	}

	c.Subscribe("delivery", func(s string) {
		fmt.Printf("Received a message: %s\n", s)
	})

	c.Subscribe("receipt", func(s string) {
		fmt.Printf("Received a message: %s\n", s)
	})

	env.NATS = c

	// Initialize factors
	env.Factors = make(map[string]factor.Factor)
	if flags.YubiCloudID != "" {
		yubicloud, err := factor.NewYubiCloud(flags.YubiCloudID, flags.YubiCloudKey)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Unable to initiate YubiCloud")
		}
		env.Factors[yubicloud.Type()] = yubicloud
	}

	authenticator := factor.NewAuthenticator(6)
	env.Factors[authenticator.Type()] = authenticator

	// Create a new goji mux
	mux := web.New()

	// Include the most basic middlewares:
	//  - RequestID assigns an unique ID for each request in order to identify errors.
	//  - Glogrus logs each request
	//  - Recoverer prevents panics from crashing the API
	//  - AutomaticOptions automatically responds to OPTIONS requests
	mux.Use(middleware.RequestID)
	mux.Use(glogrus.NewGlogrus(log, "api"))
	mux.Use(middleware.Recoverer)
	mux.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler)
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

	// Tokens
	auth.Get("/tokens", routes.TokensGet)
	mux.Post("/tokens", routes.TokensCreate)
	auth.Delete("/tokens", routes.TokensDelete)
	auth.Delete("/tokens/:id", routes.TokensDelete)

	// Threads
	auth.Get("/threads", routes.ThreadsList)
	auth.Get("/threads/:id", routes.ThreadsGet)
	auth.Put("/threads/:id", routes.ThreadsUpdate)

	// Emails
	auth.Get("/emails", routes.EmailsList)
	auth.Post("/emails", routes.EmailsCreate)
	auth.Get("/emails/:id", routes.EmailsGet)
	auth.Put("/emails/:id", routes.EmailsUpdate)
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

	// WebSockets handler
	ws, err := socketio.NewServer(nil)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to create a socket.io server")
	}
	ws.On("connection", func(so socketio.Socket) {
		env.Log.WithFields(logrus.Fields{
			"id": so.Id(),
		}).Info("New WebSockets connection")

		so.On("request", func(id string, method string, path string, data string, headers map[string]string) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(method, "http://api.lavaboom.io"+path, strings.NewReader(data))
			if err != nil {
				so.Emit("error", err.Error())
				return
			}

			for key, value := range headers {
				r.Header.Set(key, value)
			}

			mux.ServeHTTP(w, r)

			resp, err := http.ReadResponse(bufio.NewReader(w.Body), r)
			if err != nil {
				so.Emit("error", err.Error())
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				so.Emit("error", err.Error())
				return
			}

			so.Emit("response", id, resp.StatusCode, resp.Header, body)
		})
	})

	mux.Handle("/ws/*", sockjs.NewHandler("/ws", sockjs.DefaultOptions, func(session sockjs.Session) {
		// A new goroutine seems to be spawned for each new session
		for {
			// Read a message from the input
			msg, err := session.Recv()
			if err != nil {
				env.Log.WithFields(logrus.Fields{
					"id":    session.ID(),
					"error": err.Error(),
				}).Warn("Error while reading from a WebSocket")
				break
			}

			// Decode the message
			var input struct {
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

			// Perform the request
			w := httptest.NewRecorder()
			r, err := http.NewRequest(input.Method, "http://api.lavaboom.io"+input.Path, strings.NewReader(input.Body))
			if err != nil {
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

			r.RequestURI = input.Path

			for key, value := range input.Headers {
				r.Header.Set(key, value)
			}

			mux.ServeHTTP(w, r)

			// Return the final response
			result, _ := json.Marshal(map[string]interface{}{
				"type":   "response",
				"id":     input.ID,
				"status": w.Code,
				"header": w.HeaderMap,
				"body":   w.Body.String(),
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
	}))

	// Merge the muxes
	mux.Handle("/*", auth)

	// Compile the routes
	mux.Compile()

	return mux
}
