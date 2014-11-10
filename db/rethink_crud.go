package db

import (
	rethink "github.com/dancannon/gorethink"
)

// RethinkTable contains the most basic table functions
type RethinkTable interface {
	GetTableName() string
	GetDBName() string
}

// RethinkCreator contains a function to create new instances in the table
type RethinkCreator interface {
	Insert(data interface{}) error
}

// RethinkReader allows fetching resources from the database
type RethinkReader interface {
	Find(id string) (*rethink.Cursor, error)
	FindFetchOne(id string, value interface{}) error

	FindBy(key string, value interface{}) (*rethink.Cursor, error)
	FindByAndFetch(key string, value interface{}, results interface{}) error
	FindByAndFetchOne(key string, value interface{}, result interface{}) error

	Where(filter map[string]interface{}) (*rethink.Cursor, error)
	WhereAndFetch(filter map[string]interface{}, results interface{}) error
	WhereAndFetchOne(filter map[string]interface{}, result interface{}) error

	FindByIndex(index string, values ...interface{}) (*rethink.Cursor, error)
	FindByIndexFetch(results interface{}, index string, values ...interface{}) error
	FindByIndexFetchOne(result interface{}, index string, values ...interface{}) error
}

// RethinkUpdater allows updating existing resources in the database
type RethinkUpdater interface {
	Update(data interface{}) error
	UpdateID(id string, data interface{}) error
}

// RethinkDeleter allows deleting resources from the database
type RethinkDeleter interface {
	Delete(pred interface{}) error
	DeleteID(id string) error
}

// RethinkCRUD is the interface that every table should implement
type RethinkCRUD interface {
	RethinkCreator
	RethinkReader
	RethinkUpdater
	RethinkDeleter
	RethinkTable
}
