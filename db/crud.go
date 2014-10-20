package db

import (
	"errors"
	"log"

	r "github.com/dancannon/gorethink"
)

// TODO: throw custom errors

// Insert performs an insert operation for any map[string]interface{} or struct
func Insert(table string, data interface{}) error {
	return insertHelper(table, data, "error")
}

// Update performs an update operation for any map[string]interface{} or struct
func Update(table string, data interface{}) error {
	return insertHelper(table, data, "update")
}

// Delete deletes a database item based on id
func Delete(table string, id string) error {
	_, err := r.Table(table).Get(id).Delete().RunWrite(config.Session)
	if err != nil {
		log.Fatalf("Couldn't delete [%s] in table [%s]\n", id, table)
	}
	return nil
}

// Get fetches a database object with a specific id
func Get(table string, id string) (*r.Cursor, error) {
	if response, err := r.Table(table).Get(id).Run(config.Session); err == nil {
		return response, nil
	}
	return nil, errors.New("Item not found")
}

// GetAll fetches all items in the table that satisfy item[index] == value
// TODO: find out how to match on nested keys
func GetAll(table string, index string, value interface{}) (*r.Cursor, error) {
	if response, err := r.Table(table).GetAllByIndex(index, value).Run(config.Session); err == nil {
		log.Println("db.GetAll", response)
		return response, nil
	}
	return nil, errors.New("Not found")
}

// GetByID is an alias for Get
var GetByID = Get

// GetByIndex is an alias for GetAll
var GetByIndex = GetAll

// Remove is an alias for Delete
var Remove = Delete

// Rm is an alias for Delete
var Rm = Delete

// insertHelper adds an interface{} to the database. Helper func for db.Insert and db.Update
func insertHelper(table string, data interface{}, conflictResolution string) error {
	// TODO check out the RunWrite result, conflict errors are reported there
	_, err := r.Table(table).Insert(data, r.InsertOpts{Conflict: conflictResolution}).RunWrite(config.Session)
	if err != nil {
		log.Fatalln("Database insert operation failed. Data:\n", data)
	}
	return nil
}
