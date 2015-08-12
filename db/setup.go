package db

import (
	r "github.com/dancannon/gorethink"
)

// List of names of databases
var databaseNames = []string{
	"prod",
	"staging",
	"dev",
	"test",
}

// Setup configures the RethinkDB server
func Setup(opts r.ConnectOpts) error {
	// Initialize a new setup connection
	ss, err := r.Connect(opts)
	if err != nil {
		return err
	}

	// Create databases
	for _, d := range databaseNames {
		r.DBCreate(d).Exec(ss)

		r.DB(d).TableCreate("accounts").Exec(ss)
		r.DB(d).Table("accounts").IndexCreate("name").Exec(ss)
		r.DB(d).Table("accounts").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("accounts").IndexCreate("date_modified").Exec(ss)
		r.DB(d).Table("accounts").IndexCreate("alt_email").Exec(ss)
		r.DB(d).Table("accounts").IndexCreate("type").Exec(ss)
		r.DB(d).Table("accounts").IndexCreate("status").Exec(ss)

		r.DB(d).TableCreate("addresses").Exec(ss)
		r.DB(d).Table("addresses").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("addresses").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("addresses").IndexCreate("date_modified").Exec(ss)

		r.DB(d).TableCreate("contacts").Exec(ss)
		r.DB(d).Table("contacts").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("contacts").IndexCreate("name").Exec(ss)
		r.DB(d).Table("contacts").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("contacts").IndexCreate("date_modified").Exec(ss)

		r.DB(d).TableCreate("emails").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("date_modified").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("thread").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("kind").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("from").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("message_id").Exec(ss)
		r.DB(d).Table("emails").IndexCreate("to", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.DB(d).Table("emails").IndexCreate("cc", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.DB(d).Table("emails").IndexCreate("bcc", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.DB(d).Table("emails").IndexCreateFunc("messageIDOwner", func(row r.Term) interface{} {
			return []interface{}{
				row.Field("message_id"),
				row.Field("owner"),
			}
		}).Exec(ss)
		r.DB(d).Table("emails").IndexCreateFunc("threadStatus", func(row r.Term) interface{} {
			return []interface{}{
				row.Field("thread"),
				row.Field("status"),
			}
		}).Exec(ss)

		r.DB(d).TableCreate("files").Exec(ss)
		r.DB(d).Table("files").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("files").IndexCreate("name").Exec(ss)
		r.DB(d).Table("files").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("files").IndexCreate("date_modified").Exec(ss)

		r.DB(d).TableCreate("keys").Exec(ss)
		r.DB(d).Table("keys").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("keys").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("keys").IndexCreate("date_modified").Exec(ss)
		r.DB(d).Table("keys").IndexCreate("key_id").Exec(ss)

		r.DB(d).TableCreate("labels").Exec(ss)
		r.DB(d).Table("labels").IndexCreate("name").Exec(ss)
		r.DB(d).Table("labels").IndexCreate("builtin").Exec(ss)
		r.DB(d).Table("labels").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("labels").IndexCreateFunc("nameOwnerBuiltin", func(row r.Term) interface{} {
			return []interface{}{
				row.Field("name"),
				row.Field("owner"),
				row.Field("builtin"),
			}
		}).Exec(ss)

		r.DB(d).TableCreate("threads").Exec(ss)
		r.DB(d).Table("threads").IndexCreate("name").Exec(ss)
		r.DB(d).Table("threads").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("threads").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("threads").IndexCreate("date_modified").Exec(ss)
		r.DB(d).Table("threads").IndexCreate("emails", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.DB(d).Table("threads").IndexCreate("labels", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.DB(d).Table("threads").IndexCreate("members", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.DB(d).Table("threads").IndexCreate("subject_hash").Exec(ss)
		r.DB(d).Table("threads").IndexCreate("secure").Exec(ss)
		r.DB(d).Table("threads").IndexCreateFunc("subjectOwner", func(row r.Term) interface{} {
			return []interface{}{
				row.Field("subject_hash"),
				row.Field("owner"),
			}
		}).Exec(ss)

		r.DB(d).TableCreate("tokens").Exec(ss)
		r.DB(d).Table("tokens").IndexCreate("name").Exec(ss)
		r.DB(d).Table("tokens").IndexCreate("owner").Exec(ss)
		r.DB(d).Table("tokens").IndexCreate("date_created").Exec(ss)
		r.DB(d).Table("tokens").IndexCreate("date_modified").Exec(ss)
		r.DB(d).Table("tokens").IndexCreate("type").Exec(ss)
		r.DB(d).Table("tokens").IndexCreate("expiry_date").Exec(ss)

		r.DB(d).TableCreate("webhooks").Exec(ss)
		r.DB(d).Table("webhooks").IndexCreate("target").Exec(ss)
		r.DB(d).Table("webhooks").IndexCreate("type").Exec(ss)
		r.DB(d).Table("webhooks").IndexCreateFunc("targetType", func(row r.Term) interface{} {
			return []interface{}{
				row.Field("target"),
				row.Field("type"),
			}
		}).Exec(ss)
	}

	return ss.Close()
}
