package db

import (
	"time"

	"github.com/dancannon/gorethink"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/models"
)

type LabelsTable struct {
	RethinkCRUD
	Cache   cache.Cache
	Expires time.Duration
}

func (l *LabelsTable) Insert(data interface{}) error {
	if err := l.RethinkCRUD.Insert(data); err != nil {
		return err
	}

	label, ok := data.(*models.Token)
	if !ok {
		return nil
	}

	return l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+label.ID, label, l.Expires)
}

// Update clears all updated keys
func (l *LabelsTable) Update(data interface{}) error {
	if err := l.RethinkCRUD.Update(data); err != nil {
		return err
	}

	return l.Cache.DeleteMask(l.RethinkCRUD.GetTableName() + ":*")
}

// UpdateID updates the specified label and updates cache
func (l *LabelsTable) UpdateID(id string, data interface{}) error {
	if err := l.RethinkCRUD.UpdateID(id, data); err != nil {
		return err
	}

	label, err := l.GetLabel(id)
	if err != nil {
		return err
	}

	return l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+id, label, l.Expires)
}

// Delete removes from db and cache using filter
func (l *LabelsTable) Delete(cond interface{}) error {
	result, err := l.GetTable().Filter(cond).Delete(gorethink.DeleteOpts{
		ReturnChanges: true,
	}).RunWrite(l.GetSession())
	if err != nil {
		return err
	}

	var ids []interface{}
	for _, change := range result.Changes {
		ids = append(ids, l.RethinkCRUD.GetTableName()+":"+change.OldValue.(map[string]interface{})["id"].(string))
	}

	return l.Cache.DeleteMulti(ids...)
}

// DeleteID removes from db and cache using id query
func (l *LabelsTable) DeleteID(id string) error {
	if err := l.RethinkCRUD.DeleteID(id); err != nil {
		return err
	}

	return l.Cache.Delete(l.RethinkCRUD.GetTableName() + ":" + id)
}

// FindFetchOne tries cache and then tries using DefaultCRUD's fetch operation
func (l *LabelsTable) FindFetchOne(id string, value interface{}) error {
	if err := l.Cache.Get(id, value); err == nil {
		return nil
	}

	err := l.RethinkCRUD.FindFetchOne(id, value)
	if err != nil {
		return err
	}

	return l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+id, value, l.Expires)
}

func (l *LabelsTable) GetLabel(id string) (*models.Label, error) {
	var result models.Label

	if err := l.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
