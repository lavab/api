package db

import (
	"time"

	//"github.com/dancannon/gorethink"

	//"github.com/lavab/api/cache"
	"github.com/lavab/api/models"
)

type LabelsTable struct {
	RethinkCRUD
	Emails *EmailsTable
	//Cache   cache.Cache
	Expires time.Duration
}

func (l *LabelsTable) Insert(data interface{}) error {
	if err := l.RethinkCRUD.Insert(data); err != nil {
		return err
	}

	return nil

	//label, ok := data.(*models.Token)
	//if !ok {
	//	return nil
	//}

	//return l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+label.ID, label, l.Expires)
}

// Update clears all updated keys
func (l *LabelsTable) Update(data interface{}) error {
	if err := l.RethinkCRUD.Update(data); err != nil {
		return err
	}

	return nil
	//return l.Cache.DeleteMask(l.RethinkCRUD.GetTableName() + ":*")
}

// UpdateID updates the specified label and updates cache
func (l *LabelsTable) UpdateID(id string, data interface{}) error {
	if err := l.RethinkCRUD.UpdateID(id, data); err != nil {
		return err
	}

	return nil

	//label, err := l.GetLabel(id)
	//if err != nil {
	//	return err
	//}

	//return l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+id, label, l.Expires)
}

// Delete removes from db and cache using filter
func (l *LabelsTable) Delete(cond interface{}) error {
	/*result, err := l.GetTable().Filter(cond).Delete(gorethink.DeleteOpts{
		ReturnChanges: true,
	}).RunWrite(l.GetSession())*/
	_, err := l.GetTable().Filter(cond).Delete().RunWrite(l.GetSession())
	if err != nil {
		return err
	}

	return nil

	/*var ids []interface{}
	for _, change := range result.Changes {
		ids = append(ids, l.RethinkCRUD.GetTableName()+":"+change.OldValue.(map[string]interface{})["id"].(string))
	}

	return l.Cache.DeleteMulti(ids...)*/
}

// DeleteID removes from db and cache using id query
func (l *LabelsTable) DeleteID(id string) error {
	/*label, err := l.GetLabel(id)
	if err != nil {
		return err
	}*/

	if err := l.RethinkCRUD.DeleteID(l.RethinkCRUD.GetTableName() + ":" + id); err != nil {
		return err
	}

	/*l.Cache.Delete(l.RethinkCRUD.GetTableName() + ":" + id)
	l.Cache.Delete(l.RethinkCRUD.GetTableName() + ":owner:" + label.Owner)*/

	return nil
}

func (l *LabelsTable) GetLabel(id string) (*models.Label, error) {
	var result models.Label

	/*if err := l.Cache.Get(l.RethinkCRUD.GetTableName()+":"+id, &result); err == nil {
		return &result, nil
	}*/

	if err := l.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	/*err := l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+id, result, l.Expires)
	if err != nil {
		return nil, err
	}*/

	return &result, nil
}

// GetOwnedBy returns all labels owned by id
func (l *LabelsTable) GetOwnedBy(id string) ([]*models.Label, error) {
	var result []*models.Label

	/*if err := l.Cache.Get(l.RethinkCRUD.GetTableName()+":owner:"+id, &result); err == nil {
		return result, nil
	}*/

	err := l.WhereAndFetch(map[string]interface{}{
		"owner": id,
	}, &result)
	if err != nil {
		return nil, err
	}

	/*err = l.Cache.Set(l.RethinkCRUD.GetTableName()+":owner:"+id, result, l.Expires)
	if err != nil {
		return nil, err
	}*/

	return result, nil
}

func (l *LabelsTable) GetLabelByNameAndOwner(owner string, name string) (*models.Label, error) {
	var result models.Label

	if err := l.WhereAndFetchOne(map[string]interface{}{
		"name":  name,
		"owner": owner,
	}, &result); err != nil {
		return nil, err
	}

	/*err := l.Cache.Set(l.RethinkCRUD.GetTableName()+":"+id, result, l.Expires)
	if err != nil {
		return nil, err
	}*/

	return &result, nil
}
