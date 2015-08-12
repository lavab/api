package db

import (
	"github.com/dancannon/gorethink"

	"github.com/lavab/api/models"
)

// Emails implements the CRUD interface for tokens
type EmailsTable struct {
	RethinkCRUD
}

// GetEmail returns a token with specified name
func (e *EmailsTable) GetEmail(id string) (*models.Email, error) {
	var result models.Email

	if err := e.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOwnedBy returns all emails owned by id
func (e *EmailsTable) GetOwnedBy(id string) ([]*models.Email, error) {
	var result []*models.Email

	err := e.WhereAndFetch(map[string]interface{}{
		"owner": id,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOwnedBy deletes all emails owned by id
func (e *EmailsTable) DeleteOwnedBy(id string) error {
	return e.Delete(map[string]interface{}{
		"owner": id,
	})
}

func (e *EmailsTable) CountOwnedBy(id string) (int, error) {
	return e.FindByAndCount("owner", id)
}

func (e *EmailsTable) List(
	owner string,
	sort []string,
	offset int,
	limit int,
	thread string,
) ([]*models.Email, error) {

	filter := map[string]interface{}{}

	if owner != "" {
		filter["owner"] = owner
	}

	if thread != "" {
		filter["thread"] = thread
	}

	term := e.GetTable().Filter(filter).Filter(gorethink.Not(gorethink.Row.Field("status").Eq(gorethink.Expr("queued"))))

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
	cursor, err := term.Run(e.GetSession())
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	// Fetch the cursor
	var resp []*models.Email
	err = cursor.All(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *EmailsTable) GetByThread(thread string) ([]*models.Email, error) {
	var result []*models.Email

	err := e.FindByAndFetch("thread", thread, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *EmailsTable) DeleteByThread(id string) error {
	return e.Delete(map[string]interface{}{
		"thread": id,
	})
}

func (e *EmailsTable) GetThreadManifest(thread string) (string, error) {
	cursor, err := e.GetTable().
		GetAllByIndex("thread", thread).
		OrderBy("date_created").
		Limit(1).
		Pluck("manifest").
		Field("manifest").
		Run(e.GetSession())
	if err != nil {
		return "", err
	}

	var manifest string
	err = cursor.One(&manifest)
	if err != nil {
		return "", err
	}

	return manifest, nil
}
