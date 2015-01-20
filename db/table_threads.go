package db

import (
	"github.com/dancannon/gorethink"

	"github.com/lavab/api/models"
)

type ThreadsTable struct {
	RethinkCRUD
}

func (t *ThreadsTable) GetThread(id string) (*models.Thread, error) {
	var result models.Thread

	if err := t.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *ThreadsTable) GetOwnedBy(id string) ([]*models.Thread, error) {
	var result []*models.Thread

	err := t.WhereAndFetch(map[string]interface{}{
		"owner": id,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *ThreadsTable) DeleteOwnedBy(id string) error {
	return t.Delete(map[string]interface{}{
		"owner": id,
	})
}

func (t *ThreadsTable) CountOwnedBy(id string) (int, error) {
	return t.FindByAndCount("owner", id)
}

func (t *ThreadsTable) List(
	owner string,
	sort []string,
	offset int,
	limit int,
	label string,
) ([]*models.Thread, error) {

	var term gorethink.Term

	if owner != "" && label != "" {
		term = t.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
			return gorethink.And(
				row.Field("owner").Eq(owner),
				row.Field("labels").Contains(label),
			)
		})
	}

	if owner != "" && label == "" {
		term = t.GetTable().Filter(map[string]interface{}{
			"owner": owner,
		})
	}

	if owner == "" && label != "" {
		term = t.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
			return row.Field("labels").Contains(label)
		})
	}

	// If sort array has contents, parse them and add to the term
	if sort != nil && len(sort) > 0 {
		var conds []interface{}
		for _, cond := range sort {
			if cond[0] == '-' {
				conds = append(conds, gorethink.Desc(cond[1:]))
			} else if cond[0] == '+' || cond[0] == ' ' {
				conds = append(conds, gorethink.Asc(cond[1:]))
			} else {
				conds = append(conds, gorethink.Asc(cond))
			}
		}

		term = term.OrderBy(conds...)
	}

	// Slice the result in 3 cases
	if offset != 0 && limit == 0 {
		term = term.Skip(offset)
	}

	if offset == 0 && limit != 0 {
		term = term.Limit(limit)
	}

	if offset != 0 && limit != 0 {
		term = term.Slice(offset, offset+limit)
	}

	// Run the query
	cursor, err := term.Run(t.GetSession())
	if err != nil {
		return nil, err
	}

	// Fetch the cursor
	var resp []*models.Thread
	err = cursor.All(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *ThreadsTable) GetByLabel(label string) ([]*models.Thread, error) {
	var result []*models.Thread

	cursor, err := t.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
		return row.Field("labels").Contains(label)
	}).GetAll().Run(t.GetSession())
	if err != nil {
		return nil, err
	}

	err = cursor.All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *ThreadsTable) CountByLabel(label string) (int, error) {
	var result int

	cursor, err := t.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
		return row.Field("labels").Contains(label)
	}).Count().Run(t.GetSession())
	if err != nil {
		return 0, err
	}

	err = cursor.One(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (t *ThreadsTable) CountByLabelUnread(label string) (int, error) {
	var result int

	cursor, err := t.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
		return gorethink.And(
			row.Field("labels").Contains(label),
			row.Field("is_read").Eq(false),
		)
	}).Count().Run(t.GetSession())
	if err != nil {
		return 0, err
	}

	err = cursor.One(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
