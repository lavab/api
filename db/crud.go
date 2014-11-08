package db

import (
	r "github.com/dancannon/gorethink"
)

type RethinkTable interface {
	GetTableName() string
	GetDbName() string
}

type RethinkCreater interface {
	Insert(data interface{}) error
}

type RethinkReader interface {
	Find(id string) (*r.Cursor, error)
	FindFetchOne(id string, value interface{}) error

	FindBy(key string, value interface{}) (*r.Cursor, error)
	FindByAndFetch(key string, value interface{}, results interface{}) error
	FindByAndFetchOne(key string, value interface{}, result interface{}) error

	Where(filter map[string]interface{}) (*r.Cursor, error)
	WhereAndFetch(filter map[string]interface{}, results interface{}) error
	WhereAndFetchOne(filter map[string]interface{}, result interface{}) error

	FindByIndex(index string, values ...interface{}) (*r.Cursor, error)
	FindByIndexFetch(results interface{}, index string, values ...interface{}) error
	FindByIndexFetchOne(result interface{}, index string, values ...interface{}) error
}

type RethinkUpdater interface {
	Update(data interface{}) error
	UpdateId(id string, data interface{}) error
}

type RethinkDeleter interface {
	Delete(pred interface{}) error
	DeleteId(id string) error
}

//The interface that all tables should implement
type RethinkCrud interface {
	RethinkCreater
	RethinkReader
	RethinkUpdater
	RethinkDeleter
	RethinkTable
}

//The default impementation that should be embedded
type RethinkCrudImpl struct {
	table string
	db    string
}

func NewCrudTable(db, table string) *RethinkCrudImpl {
	return &RethinkCrudImpl{
		db:    db,
		table: table,
	}
}

//The RethinkTable implementation
func (rc *RethinkCrudImpl) GetTableName() string {
	return rc.table
}

func (rc *RethinkCrudImpl) GetDbName() string {
	return rc.db
}

//Gets the current table as a Rethink Term
func (rc *RethinkCrudImpl) GetTable() r.Term {
	return r.Table(rc.table)
}

//inserts a document to the database
func (rc *RethinkCrudImpl) Insert(data interface{}) error {
	_, err := rc.GetTable().Insert(data).RunWrite(config.Session)
	if err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}

//Updates according to the specified options in data
func (rc *RethinkCrudImpl) Update(data interface{}) error {
	_, err := rc.GetTable().Update(data).RunWrite(config.Session)
	if err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}

//Updates a specific id in database, with options from data
func (rc *RethinkCrudImpl) UpdateId(id string, data interface{}) error {
	_, err := rc.GetTable().Get(id).Update(data).RunWrite(config.Session)
	if err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}

//Deletes the documents that are in pred argument
func (rc *RethinkCrudImpl) Delete(pred interface{}) error {
	_, err := rc.GetTable().Filter(pred).Delete().RunWrite(config.Session)
	if err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}

//Deletes a given id
func (rc *RethinkCrudImpl) DeleteId(id string) error {
	_, err := rc.GetTable().Get(id).Delete().RunWrite(config.Session)
	if err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}

//Finds a given object from db, it does not perform any fetching
func (rc *RethinkCrudImpl) Find(id string) (*r.Cursor, error) {
	cursor, err := rc.GetTable().Get(id).Run(config.Session)
	if err != nil {
		return nil, NewDbErr(rc, err)
	}

	return cursor, nil
}

//Fetches the specified object from db and fills the value with it
func (rc *RethinkCrudImpl) FindFetchOne(id string, value interface{}) error {
	cursor, err := rc.Find(id)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.One(value); err != nil {
		return NewDbErr(rc, err)
	}

	//we have success here
	return nil
}

//FindBy is for looking up for key=value situations
func (rc *RethinkCrudImpl) FindBy(key string, value interface{}) (*r.Cursor, error) {
	filterMap := map[string]interface{}{
		key: value,
	}
	cursor, err := rc.GetTable().Filter(filterMap).Run(config.Session)
	if err != nil {
		return nil, NewDbErr(rc, err)
	}

	return cursor, nil
}

//FindBy is for looking up for key=value situations with fetch all
func (rc *RethinkCrudImpl) FindByAndFetch(key string, value interface{}, results interface{}) error {

	cursor, err := rc.FindBy(key, value)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.All(results); err != nil {
		return NewDbErr(rc, err)
	}

	//we have success here
	return nil
}

//Fetches the specified object from db and fills the value with it
func (rc *RethinkCrudImpl) FindByAndFetchOne(key string, value interface{}, result interface{}) error {

	cursor, err := rc.FindBy(key, value)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.One(result); err != nil {
		return NewDbErr(rc, err)
	}

	//we have success here
	return nil
}

//Where is for asking for more than one field from database, useful for passing a few
//pairs for AND querying
func (rc *RethinkCrudImpl) Where(filter map[string]interface{}) (*r.Cursor, error) {
	cursor, err := rc.GetTable().Filter(filter).Run(config.Session)
	if err != nil {
		return nil, NewDbErr(rc, err)
	}

	return cursor, nil
}

//Where with fetch all
func (rc *RethinkCrudImpl) WhereAndFetch(filter map[string]interface{}, results interface{}) error {
	cursor, err := rc.Where(filter)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.All(results); err != nil {
		return NewDbErr(rc, err)
	}

	//we have success here
	return nil
}

//Where with fetch all
func (rc *RethinkCrudImpl) WhereAndFetchOne(filter map[string]interface{}, result interface{}) error {
	cursor, err := rc.Where(filter)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.One(result); err != nil {
		return NewDbErr(rc, err)
	}

	//we have success here
	return nil
}

//GetAll fetches all items in the table that satisfy item[index] == value
func (rc *RethinkCrudImpl) FindByIndex(index string, values ...interface{}) (*r.Cursor, error) {
	cursor, err := rc.GetTable().GetAllByIndex(index, values...).Run(config.Session)
	if err != nil {
		return nil, NewDbErr(rc, err)
	}

	return cursor, nil
}

func (rc *RethinkCrudImpl) FindByIndexFetch(results interface{}, index string, values ...interface{}) error {
	cursor, err := rc.FindByIndex(index, values...)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.All(results); err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}

func (rc *RethinkCrudImpl) FindByIndexFetchOne(result interface{}, index string, values ...interface{}) error {
	cursor, err := rc.FindByIndex(index, values...)
	if err != nil {
		return err
	}

	//now fetch the item from  database
	if err := cursor.One(result); err != nil {
		return NewDbErr(rc, err)
	}

	return nil
}
