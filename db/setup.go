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
		r.DbCreate(d).Exec(ss)

		r.Db(d).TableCreate("accounts").Exec(ss)
		r.Db(d).Table("accounts").IndexCreate("name").Exec(ss)
		r.Db(d).Table("accounts").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("accounts").IndexCreate("date_modified").Exec(ss)
		r.Db(d).Table("accounts").IndexCreate("alt_email").Exec(ss)
		r.Db(d).Table("accounts").IndexCreate("type").Exec(ss)
		r.Db(d).Table("accounts").IndexCreate("status").Exec(ss)

		r.Db(d).TableCreate("contacts").Exec(ss)
		r.Db(d).Table("contacts").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("contacts").IndexCreate("name").Exec(ss)
		r.Db(d).Table("contacts").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("contacts").IndexCreate("date_modified").Exec(ss)

		r.Db(d).TableCreate("emails").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("date_modified").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("thread").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("kind").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("from").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("message_id").Exec(ss)
		r.Db(d).Table("emails").IndexCreate("to", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.Db(d).Table("emails").IndexCreate("cc", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.Db(d).Table("emails").IndexCreate("bcc", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.Db(d).Table("emails").IndexCreateFunc("messageIDOwner", func(row r.Term) r.Term {
			return r.Expr([]interface{}{
				row.Field("message_id"),
				row.Field("owner"),
			})
		}).Exec(ss)
		r.Db(d).Table("emails").IndexCreateFunc("threadStatus", func(row r.Term) r.Term {
			return r.Expr([]interface{}{
				row.Field("thread"),
				row.Field("status"),
			})
		})

		r.Db(d).TableCreate("files").Exec(ss)
		r.Db(d).Table("files").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("files").IndexCreate("name").Exec(ss)
		r.Db(d).Table("files").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("files").IndexCreate("date_modified").Exec(ss)

		r.Db(d).TableCreate("keys").Exec(ss)
		r.Db(d).Table("keys").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("keys").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("keys").IndexCreate("date_modified").Exec(ss)
		r.Db(d).Table("keys").IndexCreate("key_id").Exec(ss)

		r.Db(d).TableCreate("labels").Exec(ss)
		r.Db(d).Table("labels").IndexCreate("name").Exec(ss)
		r.Db(d).Table("labels").IndexCreate("builtin").Exec(ss)
		r.Db(d).Table("labels").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("labels").IndexCreateFunc("nameOwnerBuiltin", func(row r.Term) r.Term {
			return r.Expr([]interface{}{
				row.Field("name"),
				row.Field("owner"),
				row.Field("builtin"),
			})
		}).Exec(ss)

		r.Db(d).TableCreate("threads").Exec(ss)
		r.Db(d).Table("threads").IndexCreate("name").Exec(ss)
		r.Db(d).Table("threads").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("threads").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("threads").IndexCreate("date_modified").Exec(ss)
		r.Db(d).Table("threads").IndexCreate("emails", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.Db(d).Table("threads").IndexCreate("labels", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.Db(d).Table("threads").IndexCreate("members", r.IndexCreateOpts{Multi: true}).Exec(ss)
		r.Db(d).Table("threads").IndexCreate("subject_hash").Exec(ss)
		r.Db(d).Table("threads").IndexCreate("secure").Exec(ss)
		r.Db(d).Table("threads").IndexCreateFunc("subjectOwner", func(row r.Term) r.Term {
			return r.Expr([]interface{}{
				row.Field("subject_hash"),
				row.Field("owner"),
			})
		})

		r.Db(d).TableCreate("tokens").Exec(ss)
		r.Db(d).Table("tokens").IndexCreate("name").Exec(ss)
		r.Db(d).Table("tokens").IndexCreate("owner").Exec(ss)
		r.Db(d).Table("tokens").IndexCreate("date_created").Exec(ss)
		r.Db(d).Table("tokens").IndexCreate("date_modified").Exec(ss)
		r.Db(d).Table("tokens").IndexCreate("type").Exec(ss)
		r.Db(d).Table("tokens").IndexCreate("expiry_date").Exec(ss)
	}

	return ss.Close()
}
