package db

import (
	"github.com/lavab/api/models"
)

// Contacts implements the CRUD interface for tokens
type ContactsTable struct {
	RethinkCRUD
}

// GetContact returns a token with specified name
func (c *ContactsTable) GetContact(id string) (*models.Contact, error) {
	var result models.Contact

	if err := c.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOwnedBy returns all contacts owned by id
func (c *ContactsTable) GetOwnedBy(id string) ([]*models.Contact, error) {
	var result []*models.Contact

	err := c.WhereAndFetch(map[string]interface{}{
		"owner": id,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOwnedBy deletes all contacts owned by id
func (c *ContactsTable) DeleteOwnedBy(id string) error {
	return c.Delete(map[string]interface{}{
		"owner": id,
	})
}
