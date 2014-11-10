package db

import (
	"fmt"
	"log"
	"os"
	"time"

	r "github.com/dancannon/gorethink"
)

const (
	TABLE_SESSIONS = "sessions"
	TABLE_USERS    = "users"
	TABLE_EMAILS   = "emails"
	TABLE_DRAFTS   = "drafts"
	TABLE_CONTACTS = "contacts"
	TABLE_THREADS  = "threads"
	TABLE_LABELS   = "labels"
	TABLE_KEYS     = "keys"
)

var config struct {
	Session *r.Session
	Url     string
	AuthKey string
	Db      string
}

var CurrentConfig = config

var dbs = []string{
	"prod",
	"staging",
	"dev",
}

var tablesAndIndexes = map[string][]string{
	TABLE_SESSIONS: []string{"user", "user_id"},
	TABLE_USERS:    []string{"name"},
	TABLE_EMAILS:   []string{"user_id"},
	TABLE_DRAFTS:   []string{"user_id"},
	TABLE_CONTACTS: []string{},
	TABLE_THREADS:  []string{"user_id"},
	TABLE_LABELS:   []string{},
	TABLE_KEYS:     []string{},
}

func init() {
	config.Url = "localhost:28015"
	config.AuthKey = ""
	config.Db = "dev"

	if tmp := os.Getenv("RETHINKDB_URL"); tmp != "" {
		config.Url = tmp
	} else if tmp := os.Getenv("RETHINKDB_PORT_28015_TCP_ADDR"); tmp != "" {
		config.Url = fmt.Sprintf("%s:28015", tmp)
	} else {
		log.Printf("No database URL specified, using %s.\n", config.Url)
	}
	if tmp := os.Getenv("RETHINKDB_AUTHKEY"); tmp != "" {
		config.AuthKey = tmp
	} else {
		log.Fatalln("Variable RETHINKDB_AUTHKEY not set.")
	}
	if tmp := os.Getenv("API_ENV"); tmp != "" {
		// TODO add check that tmp is in dbs
		config.Db = tmp
	} else {
		log.Printf("No database specified, using %s.\n", config.Db)
	}

	// Initialise databases, tables, and indexes. This might take a while if they don't exist
	setupSession, err := r.Connect(r.ConnectOpts{
		Address:     config.Url,
		AuthKey:     config.AuthKey,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
	}
	log.Println("Creating dbs \\ tables \\ indexes")
	for _, d := range dbs {
		log.Println(d)
		r.DbCreate(d).Run(setupSession)
		for t, indexes := range tablesAndIndexes {
			log.Println("›  ", t)
			r.Db(d).TableCreate(t).RunWrite(setupSession)
			for _, index := range indexes {
				log.Println("›  ›  ", index)
				r.Db(d).Table(t).IndexCreate(index).Exec(setupSession)
			}
		}
	}
	setupSession.Close()

	// Setting up the main session
	config.Session, err = r.Connect(r.ConnectOpts{
		Address:     config.Url,
		AuthKey:     config.AuthKey,
		Database:    config.Db,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})
}
