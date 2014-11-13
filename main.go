package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
	"github.com/lavab/glogrus"
	"github.com/namsral/flag"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/lavab/api/db"
	"github.com/lavab/api/env"
	"github.com/lavab/api/routes"
)

// TODO: "Middleware that implements a few quick security wins"
// 		 https://github.com/unrolled/secure

var (
	// General flags
	bindAddress      = flag.String("bind", ":5000", "Network address used to bind")
	apiVersion       = flag.String("version", "v0", "Shown API version")
	logFormatterType = flag.String("log", "text", "Log formatter type. Either \"json\" or \"text\"")
	sessionDuration  = flag.Int("session_duration", 72, "Session duration expressed in hours")
	forceColors      = flag.Bool("force_colors", false, "Force colored prompt?")
	// Database-related flags
	rethinkdbURL = flag.String("rethinkdb_url", func() string {
		address := os.Getenv("RETHINKDB_PORT_28015_TCP_ADDR")
		if address == "" {
			address = "localhost"
		}
		return address + ":28015"
	}(), "Address of the RethinkDB database")
	rethinkdbKey      = flag.String("rethinkdb_key", os.Getenv("RETHINKDB_AUTHKEY"), "Authentication key of the RethinkDB database")
	rethinkdbDatabase = flag.String("rethinkdb_db", func() string {
		database := os.Getenv("RETHINKDB_NAME")
		if database == "" {
			database = "dev"
		}
		return database
	}(), "Database name on the RethinkDB server")
)

func main() {
	// Parse the flags
	flag.Parse()

	// Put config into the environment package
	env.Config = &env.Flags{
		BindAddress:      *bindAddress,
		APIVersion:       *apiVersion,
		LogFormatterType: *logFormatterType,
		SessionDuration:  *sessionDuration,
	}

	// Set up a new logger
	log := logrus.New()

	// Set the formatter depending on the passed flag's value
	if *logFormatterType == "text" {
		log.Formatter = &logrus.TextFormatter{
			ForceColors: *forceColors,
		}
	} else if *logFormatterType == "json" {
		log.Formatter = &logrus.JSONFormatter{}
	}

	// Pass it to the environment package
	env.Log = log

	// Set up the database
	rethinkOpts := gorethink.ConnectOpts{
		Address:     *rethinkdbURL,
		AuthKey:     *rethinkdbKey,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	}
	err := db.Setup(rethinkOpts)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to set up the database")
	}

	// Initialize the actual connection
	rethinkOpts.Database = *rethinkdbDatabase
	rethinkSession, err := gorethink.Connect(rethinkOpts)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to connect to the database")
	}

	// Put the RethinkDB session into the environment package
	env.Rethink = rethinkSession

	// Initialize the tables
	env.Accounts = &db.AccountsTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"accounts",
		),
	}
	env.Tokens = &db.TokensTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"tokens",
		),
	}
	env.Keys = &db.KeysTable{
		RethinkCRUD: db.NewCRUDTable(
			rethinkSession,
			rethinkOpts.Database,
			"keys",
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
	auth.Get("/accounts/:id/sessions", routes.AccountsSessionsList)

	// Tokens
	auth.Get("/tokens", routes.TokensGet)
	mux.Post("/tokens", routes.TokensCreate)
	auth.Delete("/tokens", routes.TokensDelete)

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

	// Make the mux handle every request
	http.Handle("/", mux)

	// Log that we're starting the server
	log.WithFields(logrus.Fields{
		"address": *bindAddress,
	}).Info("Starting the HTTP server")

	// Initialize the goroutine listening to signals passed to the app
	graceful.HandleSignals()

	// Pre-graceful shutdown event
	graceful.PreHook(func() {
		log.Info("Received a singnal, stopping the application")
	})

	// Post-shutdown event
	graceful.PostHook(func() {
		log.Info("Stopped the application")
	})

	// Listen to the passed address
	listener, err := net.Listen("tcp", *bindAddress)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"address": *bindAddress,
		}).Fatal("Cannot set up a TCP listener")
	}

	// Start the listening
	err = graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		// Don't use .Fatal! We need the code to shut down properly.
		log.Error(err)
	}

	// If code reaches this place, it means that it was forcefully closed.

	// Wait until open connections close.
	graceful.Wait()
}
