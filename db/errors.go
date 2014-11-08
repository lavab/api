package db

import (
	"fmt"
)

type DbError struct {
	err   error
	msg   string
	table RethinkTable
}

func (derr *DbError) Error() string {
	return fmt.Sprintf(
		"%s - Db: %s - Table : %s - %s ",
		derr.msg,
		derr.table.GetDbName(),
		derr.table.GetTableName(),
		derr.err)
}

func NewDbErr(t RethinkTable, err error) *DbError {
	return &DbError{
		err:   err,
		table: t,
	}
}

func NewDbErrWithMsg(t RethinkTable, err error, msg string) *DbError {
	return &DbError{
		err:   err,
		table: t,
		msg:   msg,
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
