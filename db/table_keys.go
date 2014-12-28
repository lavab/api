package db

import (
	"github.com/lavab/api/models"
)

type KeysTable struct {
	RethinkCRUD
}

func (k *KeysTable) FindByOwner(id string) ([]*models.Key, error) {
	var results []*models.Key

	if err := k.FindByAndFetch("owner", id, &results); err != nil {
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
