package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lavab/api/db"
	"github.com/lavab/api/utils"
	"github.com/stretchr/graceful"
)

// TODO: "Middleware that implements a few quick security wins"
// 		 https://github.com/unrolled/secure

const (
	cTlsFilePub  = ".tls/pub"
	cTlsFilePriv = ".tls/priv"
	cTcpPort     = 5000
	apiVersion   = "v0"
)

var config struct {
	Port         int
	PortString   string
	Host         string
	TlsAvailable bool
	MethodsJSON  string
}

func init() {
	config.Port = cTcpPort
	config.Host = ""
	config.TlsAvailable = false
	config.MethodsJSON = listRoutesString()

	if tmp := os.Getenv("API_PORT"); tmp != "" {
		tmp2, err := strconv.Atoi(tmp)
		if err != nil {
			config.Port = tmp2
		}
		log.Println("Running on non-default port", config.Port)
	}
	config.PortString = fmt.Sprintf(":%d", config.Port)

	if utils.FileExists(cTlsFilePub) && utils.FileExists(cTlsFilePriv) {
		config.TlsAvailable = true
		log.Println("Imported TLS cert/key successfully.")
	} else {
		log.Printf("TLS cert (%s) and key (%s) not found, serving plain HTTP.\n", cTlsFilePub, cTlsFilePriv)
	}

	// Set up RethinkDB
	go db.Init()
}

func main() {
	setupAndRun()
	// debug()
}

func setupAndRun() {
	r := mux.NewRouter()

	if config.TlsAvailable {
		r = r.Schemes("https").Subrouter()
	}
	if tmp := os.Getenv("API_HOST"); tmp != "" {
		r = r.Host(tmp).Subrouter()
	}

	for _, rt := range publicRoutes {
		r.HandleFunc(rt.Path, rt.HandleFunc).Methods(rt.Method)
	}

	for _, rt := range authRoutes {
		r.HandleFunc(rt.Path, AuthWrapper(rt.HandleFunc)).Methods(rt.Method)
	}

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    config.PortString,
			Handler: r,
		},
	}

	if config.TlsAvailable {
		log.Fatal(srv.ListenAndServeTLS(cTlsFilePub, cTlsFilePriv))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}

func debug() {
}
