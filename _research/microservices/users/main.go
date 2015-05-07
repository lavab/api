package main

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/lavab/goji"
	"github.com/namsral/flag"
)

var (
	configFlag        = flag.String("config", "", "config file to load")
	rethinkdbHosts    = flag.String("rethinkdb_hosts", "127.0.0.1:28015", "Addresses of the RethinkDB servers")
	rethinkdbDatabase = flag.String("rethinkdb_database", "prod", "Name of the RethinkDB database to use")
)

var (
	session *r.Session
)

func main() {
	flag.Parse()

	var err error
	session, err := r.Connect(r.ConnectOpts{
		Hosts:    strings.Split(*rethinkdbName, ","),
		Database: *rethinkdbDatabase,
	})
	if err != nil {
		log.Fatal(err)
	}

	goji.Get("/accounts", listAccounts)
	goji.Post("/accounts", createAccount)
	goji.Get("/accounts/:id", getAccount)
	goji.Put("/accounts/:id", updateAccount)
	goji.Delete("/accounts/:id", deleteAccount)

	goji.Get("/accounts/:id/tokens", listAccountTokens)
	goji.Delete("/accounts/:id/tokens", deleteAccountToken)

	goji.Get("/tokens", listTokens)
	goji.Post("/tokens", createToken)
	goji.Get("/tokens/:id", getToken)
	goji.Put("/tokens/:id", updateToken)
	goji.Delete("/tokens/:id", deleteToken)

	goji.Serve()
}
