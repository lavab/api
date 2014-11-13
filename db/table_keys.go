package db

import (
	"github.com/lavab/api/models"
)

type KeysTable struct {
	RethinkCRUD
}

func (k *KeysTable) FindByName(name string) ([]*models.Key, error) {
	var results []*models.Key

	if err := k.FindByAndFetch("owner_name", name, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (k *KeysTable) FindByFingerprint(fp string) (*models.Key, error) {
	var result models.Key

	if err := k.FindFetchOne(fp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
