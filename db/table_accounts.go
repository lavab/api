package db

import (
	"github.com/lavab/api/models"
)

// AccountsTable implements the CRUD interface for accounts
type AccountsTable struct {
	RethinkCRUD
}

// GetAccount returns an account with specified ID
func (users *AccountsTable) GetAccount(id string) (*models.Account, error) {
	var result models.Account

	if err := users.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// FindAccountByName returns an account with specified name
func (users *AccountsTable) FindAccountByName(name string) (*models.Account, error) {
	var result models.Account

	if err := users.FindByIndexFetchOne(&result, "name", name); err != nil {
		return nil, err
	}

	return &result, nil
}
