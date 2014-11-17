package db

import (
	"errors"

	"github.com/lavab/api/models"
)

// AccountsTable implements the CRUD interface for accounts
type AccountsTable struct {
	RethinkCRUD

	Tokens *TokensTable
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

func (a *AccountsTable) GetTokenOwner(token *models.Token) (*models.Account, error) {
	user, err := a.GetAccount(token.Owner)
	if err != nil {
		// Try to remove the orphaned token
		a.Tokens.DeleteID(token.ID)
		return nil, errors.New("Account disabled")
	}

	return user, nil
}
