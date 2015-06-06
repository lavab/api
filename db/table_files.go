package db

import (
	"github.com/lavab/api/models"

	"github.com/dancannon/gorethink"
)

type FilesTable struct {
	RethinkCRUD
	Emails *EmailsTable
}

func (f *FilesTable) GetFile(id string) (*models.File, error) {
	var result models.File

	if err := f.FindFetchOne(id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (f *FilesTable) GetFiles(ids ...string) ([]*models.File, error) {
	iids := make([]interface{}, len(ids))
	for i, v := range ids {
		iids[i] = v
	}

	query, err := f.GetTable().GetAll(iids...).Run(f.GetSession())
	if err != nil {
		return nil, err
	}

	var result []*models.File
	if err := query.All(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (f *FilesTable) GetOwnedBy(id string) ([]*models.File, error) {
	var result []*models.File

	err := f.WhereAndFetch(map[string]interface{}{
		"owner": id,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *FilesTable) DeleteOwnedBy(id string) error {
	return f.Delete(map[string]interface{}{
		"owner": id,
	})
}

func (f *FilesTable) GetEmailFiles(id string) ([]*models.File, error) {
	email, err := f.Emails.GetEmail(id)
	if err != nil {
		return nil, err
	}

	query, err := f.Emails.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
		return gorethink.Expr(email.Files).Contains(row.Field("id"))
	}).GetAll().Run(f.Emails.GetSession())
	if err != nil {
		return nil, err
	}

	var result []*models.File
	err = query.All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *FilesTable) CountByEmail(id string) (int, error) {
	query, err := f.GetTable().GetAllByIndex("owner", id).Count().Run(f.GetSession())
	if err != nil {
		return 0, err
	}

	var result int
	err = query.One(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (f *FilesTable) GetInEmail(owner string, email string, name string) ([]*models.File, error) {
	e, err := f.Emails.GetEmail(email)
	if err != nil {
		return nil, err
	}

	query, err := f.GetTable().Filter(func(row gorethink.Term) gorethink.Term {
		return gorethink.And(
			row.Field("owner").Eq(gorethink.Expr(owner)),
			gorethink.Expr(e.Files).Contains(row.Field("id")),
			row.Field("name").Eq(gorethink.Expr(name)),
		)
	}).Run(f.GetSession())
	if err != nil {
		return nil, err
	}

	var result []*models.File
	err = query.All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
