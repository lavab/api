package db

import (
	"github.com/lavab/api/models"
)

type AddressesTable struct {
	RethinkCRUD
}

func (a *AddressesTable) GetAddress(id string) (*models.Address, error) {
	var result models.Address
	if err := a.FindFetchOne(id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AddressesTable) GetOwnedBy(id string) ([]*models.Address, error) {
	cursor, err := a.GetTable().GetAllByIndex("owner", id).Run(a.GetSession())
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var result []*models.Address
	if err := cursor.All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AddressesTable) DeleteOwnedBy(id string) error {
	return a.GetTable().GetAllByIndex("owner", id).Delete().Exec(a.GetSession())
}
