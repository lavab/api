package db

import (
	"time"

	"github.com/dancannon/gorethink"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/models"
)

// TokensTable implements the CRUD interface for tokens
type TokensTable struct {
	RethinkCRUD
	Cache   cache.Cache
	Expires time.Duration
}

// Insert monkey-patches the DefaultCRUD method and introduces caching
func (t *TokensTable) Insert(data interface{}) error {
	if err := t.RethinkCRUD.Insert(data); err != nil {
		return err
	}

	token, ok := data.(*models.Token)
	if !ok {
		return nil
	}

	return t.Cache.Set(t.RethinkCRUD.GetTableName()+":"+token.ID, token, t.Expires)
}

// Update clears all updated keys
func (t *TokensTable) Update(data interface{}) error {
	if err := t.RethinkCRUD.Update(data); err != nil {
		return err
	}

	return t.Cache.DeleteMask(t.RethinkCRUD.GetTableName() + ":*")
}

// UpdateID updates the specified token and updates cache
func (t *TokensTable) UpdateID(id string, data interface{}) error {
	if err := t.RethinkCRUD.UpdateID(id, data); err != nil {
		return err
	}

	token, err := t.GetToken(id)
	if err != nil {
		return err
	}

	return t.Cache.Set(t.RethinkCRUD.GetTableName()+":"+id, token, t.Expires)
}

// Delete removes from db and cache using filter
func (t *TokensTable) Delete(cond interface{}) error {
	result, err := t.GetTable().Filter(cond).Delete(gorethink.DeleteOpts{
		ReturnChanges: true,
	}).RunWrite(t.GetSession())
	if err != nil {
		return err
	}

	var ids []interface{}
	for _, change := range result.Changes {
		ids = append(ids, t.RethinkCRUD.GetTableName()+":"+change.OldValue.(map[string]interface{})["id"].(string))
	}

	return t.Cache.DeleteMulti(ids...)
}

// DeleteID removes from db and cache using id query
func (t *TokensTable) DeleteID(id string) error {
	if err := t.RethinkCRUD.DeleteID(id); err != nil {
		return err
	}

	return t.Cache.Delete(t.RethinkCRUD.GetTableName() + ":" + id)
}

// FindFetchOne tries cache and then tries using DefaultCRUD's fetch operation
func (t *TokensTable) FindFetchOne(id string, value interface{}) error {
	if err := t.Cache.Get(t.RethinkCRUD.GetTableName()+":"+id, value); err == nil {
		return nil
	}

	err := t.RethinkCRUD.FindFetchOne(id, value)
	if err != nil {
		return err
	}

	return t.Cache.Set(t.RethinkCRUD.GetTableName()+":"+id, value, t.Expires)
}

// GetToken returns a token with specified name
func (t *TokensTable) GetToken(id string) (*models.Token, error) {
	var result models.Token

	if err := t.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteOwnedBy deletes all tokens owned by id
func (t *TokensTable) DeleteOwnedBy(id string) error {
	return t.Delete(map[string]interface{}{
		"owner": id,
	})
}
