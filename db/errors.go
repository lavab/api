package db

import (
	"fmt"
)

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

func NewDatabaseError(t RethinkTable, err error, message string) *DatabaseError {
	return &DatabaseError{
		err:     err,
		table:   t,
		message: message,
	}
}

type ConnectionError struct {
	error
	WrongAuthKey bool
	Unreachable  bool
}

type NotFound struct {
	error
	Database bool
	Table    bool
	Item     bool
}
