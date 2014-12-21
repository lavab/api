package db

import (
	"github.com/lavab/api/models"
)

// Emails implements the CRUD interface for tokens
type EmailsTable struct {
	RethinkCRUD
}

// GetEmail returns a token with specified name
func (c *EmailsTable) GetEmail(id string) (*models.Email, error) {
	var result models.Email

	if err := c.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOwnedBy returns all contacts owned by id
func (c *EmailsTable) GetOwnedBy(id string) ([]*models.Email, error) {
	var result []*models.Email

	err := c.WhereAndFetch(map[string]interface{}{
		"owner": id,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOwnedBy deletes all contacts owned by id
func (c *EmailsTable) DeleteOwnedBy(id string) error {
	return c.Delete(map[string]interface{}{
		"owner": id,
	})
}
