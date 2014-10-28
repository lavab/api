package db

import (
	"fmt"
	"log"
	"os"
	"time"

	r "github.com/dancannon/gorethink"
)

var config struct {
	Session *r.Session
	Url     string
	AuthKey string
	Db      string
}

var dbs = []string{
	"prod",
	"staging",
	"dev",
}

var tablesAndIndexes = map[string][]string{
	"sessions": []string{"user", "user_id"},
	"users":    []string{"name"},
	"emails":   []string{"user_id"},
	"drafts":   []string{"user_id"},
	"contacts": []string{},
	"threads":  []string{"user_id"},
	"labels":   []string{},
	"keys":     []string{},
}

func Init() {
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
