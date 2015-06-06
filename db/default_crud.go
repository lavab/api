package db

import (
	"github.com/dancannon/gorethink"
)

// Default contains the basic implementation of the gorethinkCRUD interface
type Default struct {
	table   string
	db      string
	session *gorethink.Session
}

// NewCRUDTable sets up a new Default struct
func NewCRUDTable(session *gorethink.Session, db, table string) *Default {
	return &Default{
		db:      db,
		table:   table,
		session: session,
	}
}

// GetTableName returns table's name
func (d *Default) GetTableName() string {
	return d.table
}

// GetDBName returns database's name
func (d *Default) GetDBName() string {
	return d.db
}

// GetTable returns the table as a gorethink.Term
func (d *Default) GetTable() gorethink.Term {
	return gorethink.Table(d.table)
}

// GetSession returns the current session
func (d *Default) GetSession() *gorethink.Session {
	return d.session
}

// Insert inserts a document into the database
func (d *Default) Insert(data interface{}) error {
	err := d.GetTable().Insert(data).Exec(d.session)
	if err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// Update performs an update on an existing resource according to passed data
func (d *Default) Update(data interface{}) error {
	err := d.GetTable().Update(data).Exec(d.session)
	if err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// UpdateID performs an update on an existing resource with ID that equals the id argument
func (d *Default) UpdateID(id string, data interface{}) error {
	err := d.GetTable().Get(id).Update(data).Exec(d.session)
	if err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// Delete deletes resources that match the passed filter
func (d *Default) Delete(pred interface{}) error {
	err := d.GetTable().Filter(pred).Delete().Exec(d.session)
	if err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// DeleteID deletes a resource with specified ID
func (d *Default) DeleteID(id string) error {
	err := d.GetTable().Get(id).Delete().Exec(d.session)
	if err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// Find searches for a resource in the database and then returns a cursor
func (d *Default) Find(id string) (*gorethink.Cursor, error) {
	cursor, err := d.GetTable().Get(id).Run(d.session)
	if err != nil {
		return nil, NewDatabaseError(d, err, "")
	}

	return cursor, nil
}

// FindFetchOne searches for a resource and then unmarshals the first row into value
func (d *Default) FindFetchOne(id string, value interface{}) error {
	cursor, err := d.Find(id)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.One(value); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// FindBy is an utility for fetching values if they are stored in a key-value manenr.
func (d *Default) FindBy(key string, value interface{}) (*gorethink.Cursor, error) {
	filterMap := map[string]interface{}{
		key: value,
	}

	cursor, err := d.GetTable().Filter(filterMap).Run(d.session)
	if err != nil {
		return nil, NewDatabaseError(d, err, "")
	}

	return cursor, nil
}

func (d *Default) FindByAndCount(key string, value interface{}) (int, error) {
	cursor, err := d.GetTable().Filter(map[string]interface{}{
		key: value,
	}).Count().Run(d.session)
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var count int
	if err := cursor.One(&count); err != nil {
		return 0, NewDatabaseError(d, err, "")
	}

	return count, nil
}

// FindByAndFetch retrieves a value by key and then fills results with the result.
func (d *Default) FindByAndFetch(key string, value interface{}, results interface{}) error {
	cursor, err := d.FindBy(key, value)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.All(results); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// FindByAndFetchOne retrieves a value by key and then fills result with the first row of the result
func (d *Default) FindByAndFetchOne(key string, value interface{}, result interface{}) error {
	cursor, err := d.FindBy(key, value)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.One(result); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// Where allows filtering with multiple fields
func (d *Default) Where(filter map[string]interface{}) (*gorethink.Cursor, error) {
	cursor, err := d.GetTable().Filter(filter).Run(d.session)
	if err != nil {
		return nil, NewDatabaseError(d, err, "")
	}

	return cursor, nil
}

// WhereAndFetch filters with multiple fields and then fills results with all found resources
func (d *Default) WhereAndFetch(filter map[string]interface{}, results interface{}) error {
	cursor, err := d.Where(filter)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.All(results); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// WhereAndFetchOne filters with multiple fields and then fills result with the first found resource
func (d *Default) WhereAndFetchOne(filter map[string]interface{}, result interface{}) error {
	cursor, err := d.Where(filter)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.One(result); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// FindByIndex filters all resources whose index is matching
func (d *Default) FindByIndex(index string, values ...interface{}) (*gorethink.Cursor, error) {
	cursor, err := d.GetTable().GetAllByIndex(index, values...).Run(d.session)
	if err != nil {
		return nil, NewDatabaseError(d, err, "")
	}

	return cursor, nil
}

// FindByIndexFetch filters all resources whose index is matching and fills results with all found resources
func (d *Default) FindByIndexFetch(results interface{}, index string, values ...interface{}) error {
	cursor, err := d.FindByIndex(index, values...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	//now fetch the item from  database
	if err := cursor.All(results); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}

// FindByIndexFetchOne filters all resources whose index is matching and fills result with the first one found
func (d *Default) FindByIndexFetchOne(result interface{}, index string, values ...interface{}) error {
	cursor, err := d.FindByIndex(index, values...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.One(result); err != nil {
		return NewDatabaseError(d, err, "")
	}

	return nil
}
