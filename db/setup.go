package db

import (
	"github.com/dancannon/gorethink"
)

// Publicly exported table models
var (
	Accounts *AccountsTable
	Sessions *TokensTable
)

// Indexes of tables in the database
var tableIndexes = map[string][]string{
	"tokens":       []string{"user", "user_id"},
	"accounts":     []string{"name"},
	"emails":       []string{"user_id"},
	"drafts":       []string{"user_id"},
	"contacts":     []string{},
	"threads":      []string{"user_id"},
	"labels":       []string{},
	"keys":         []string{},
	"reservations": []string{},
}

// List of names of databases
var databaseNames = []string{
	"prod",
	"staging",
	"dev",
	"test",
}

// Setup configures the RethinkDB server
func Setup(opts gorethink.ConnectOpts) error {
	// Initialize a new setup connection
	setupSession, err := gorethink.Connect(opts)
	if err != nil {
		return err
	}

	// Create databases
	for _, d := range databaseNames {
		gorethink.DbCreate(d).Run(setupSession)

		// Create tables
		for t, indexes := range tableIndexes {
			gorethink.Db(d).TableCreate(t).RunWrite(setupSession)

			// Create indexes
			for _, index := range indexes {
				gorethink.Db(d).Table(t).IndexCreate(index).Exec(setupSession)
			}
		}
	}

	return setupSession.Close()
}
