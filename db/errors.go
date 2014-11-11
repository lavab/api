package db

import (
	"fmt"
)

// DatabaseError is the wrapper for RethinkDB errors that allows passing more data with the message
type DatabaseError struct {
	err     error
	message string
	table   RethinkTable
}

func (d *DatabaseError) Error() string {
	return fmt.Sprintf(
		"%s - DB: %s - Table : %s - %s",
		d.message,
		d.table.GetDBName(),
		d.table.GetTableName(),
		d.err,
	)
}

// NewDatabaseError creates a new DatabaseError, wraps err and adds a message
func NewDatabaseError(t RethinkTable, err error, message string) *DatabaseError {
	return &DatabaseError{
		err:     err,
		table:   t,
		message: message,
	}
}
