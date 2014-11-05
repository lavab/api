package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/goji/glogrus"
	"github.com/namsral/flag"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// TODO: "Middleware that implements a few quick security wins"
// 		 https://github.com/unrolled/secure

var (
	bindAddress      = flag.String("bind", ":5000", "Network address used to bind")
	apiVersion       = flag.String("version", "v0", "Shown API version")
	logFormatterType = flag.String("log", "text", "Log formatter type. Either \"json\" or \"text\"")
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
	//  - AutomaticOptions
	mux.Use(middleware.RequestID)
	mux.Use(glogrus.NewGlogrus(log, "api"))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.AutomaticOptions)

	// Compile the routes
	mux.Compile()

	// Make the mux handle every request
	http.Handle("/", DefaultMux)

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

	// Start the listening
	err := graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		// Don't use .Fatal! We need the code to shut down properly.
		log.Error(err)
	}

	// If code reaches this place, it means that it was forcefully closed.

	// Wait until open connections close.
	graceful.Wait()
}
