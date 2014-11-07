package main

import (
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/goji/glogrus"
	"github.com/namsral/flag"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/lavab/api/env"
	"github.com/lavab/api/routes"
)

// TODO: "Middleware that implements a few quick security wins"
// 		 https://github.com/unrolled/secure

var (
	bindAddress      = flag.String("bind", ":5000", "Network address used to bind")
	apiVersion       = flag.String("version", "v0", "Shown API version")
	logFormatterType = flag.String("log", "text", "Log formatter type. Either \"json\" or \"text\"")
	sessionDuration  = flag.Int("session_duration", 72, "Session duration expressed in hours")
)

func main() {
	// Parse the flags
	flag.Parse()

	// Set up a new logger
	log := logrus.New()

	// Set the formatter depending on the passed flag's value
	if *logFormatterType == "text" {
		log.Formatter = &logrus.TextFormatter{}
	} else if *logFormatterType == "json" {
		log.Formatter = &logrus.JSONFormatter{}
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
	mux.Use(routes.AuthMiddleware)

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
	auth.Post("/tokens", routes.TokensCreate)
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
	auth.Post("/keys", routes.KeysCreate)
	mux.Get("/keys/:id", routes.KeysGet)
	auth.Post("/keys/:id/vote", routes.KeysVote)

	// Merge the muxes
	mux.Handle("/", auth)

	// Compile the routes
	mux.Compile()

	// Make the mux handle every request
	http.Handle("/", mux)

	// Set up a new environment object
	env.G = &env.Environment{
		Log: log,
		Config: &env.Config{
			BindAddress:      *bindAddress,
			APIVersion:       *apiVersion,
			LogFormatterType: *logFormatterType,
			SessionDuration:  *sessionDuration,
		},
	}

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
