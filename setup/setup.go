package setup

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/db"
	"github.com/lavab/api/env"
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

	// Merge the muxes
	mux.Handle("/*", auth)

	// Compile the routes
	mux.Compile()

	return mux
}
