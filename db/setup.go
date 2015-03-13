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
	"accounts":     []string{"name", "alt_email"},
	"contacts":     []string{"owner"},
	"emails":       []string{"owner", "label_ids", "date_created", "thread"},
	"files":        []string{"owner"},
	"keys":         []string{"owner", "key_id"},
	"labels":       []string{"owner"},
	"reservations": []string{"email", "name"},
	"threads":      []string{"owner", "subject_hash"},
	"tokens":       []string{"owner"},
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
