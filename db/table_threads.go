package db

import (
	"time"

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
	labels []string,
) ([]*models.Thread, error) {

	term := t.GetTable()

	if owner != "" {
		term = t.GetTable().GetAllByIndex("owner", owner)
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

	// Parse labels
	hasLabels := []string{}
	excLabels := []string{}
	for _, label := range labels {
		if label[0] == '-' {
			excLabels = append(excLabels, label[1:])
		} else {
			hasLabels = append(hasLabels, label)
		}
	}

	// Transform that into a term
	if len(hasLabels) > 0 || len(excLabels) > 0 {
		var hasTerm gorethink.Term
		if len(hasLabels) == 1 {
			hasTerm = gorethink.Row.Field("labels").Contains(hasLabels[0])
		} else if len(hasLabels) > 0 {
			for i, label := range hasLabels {
				if i == 0 {
					hasTerm = gorethink.Row.Field("labels").Contains(label)
				} else {
					hasTerm = hasTerm.And(gorethink.Row.Field("labels").Contains(label))
				}
			}
		}

		var excTerm gorethink.Term
		if len(excLabels) == 1 {
			excTerm = gorethink.Not(gorethink.Row.Field("labels").Contains(excLabels[0]))
		} else {
			for i, label := range excLabels {
				if i == 0 {
					excTerm = gorethink.Not(gorethink.Row.Field("labels").Contains(label))
				} else {
					excTerm = excTerm.And(gorethink.Not(gorethink.Row.Field("labels").Contains(label)))
				}
			}
		}

		// Append them into the term
		if len(hasLabels) > 0 && len(excLabels) > 0 {
			term = term.Filter(hasTerm.And(excTerm))
		} else if len(hasLabels) > 0 && len(excLabels) == 0 {
			term = term.Filter(hasTerm)
		} else if len(hasLabels) == 0 && len(excLabels) > 0 {
			term = term.Filter(excTerm)
		}
	}

	// Slice the result
	if offset != 0 || limit != 0 {
		term = term.Slice(offset, offset+limit)
	}

	// Add manifests
	term = term.Map(func(thread gorethink.Term) gorethink.Term {
		return thread.Merge(gorethink.DB(t.GetDBName()).Table("emails").Between([]interface{}{
			thread.Field("id"),
			time.Date(1990, time.January, 1, 23, 0, 0, 0, time.UTC),
		}, []interface{}{
			thread.Field("id"),
			time.Date(2090, time.January, 1, 23, 0, 0, 0, time.UTC),
		}, gorethink.BetweenOpts{
			Index: "threadAndDate",
		}).OrderBy(gorethink.OrderByOpts{Index: "threadAndDate"}).
			Nth(0).Default(map[string]interface{}{
			"interface": nil,
		}).Pluck("manifest"))
	})

	// Run the query
	cursor, err := term.Run(t.GetSession())
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

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
