package dbutils

import (
	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
)

//implements the base crud interface
type UsersTable struct {
	db.RethinkCrud
}

func (users *UsersTable) GetUser(id string) (*models.User, bool) {
	var result models.User

	if err := users.FindFetchOne(id, &result); err != nil {
		log.Println(err.Error())
		return nil, false
	}

	return &result, true

}

func (users *UsersTable) FindUserByName(username string) (*models.User, bool) {
	var result models.User

	if err := users.FindByIndexFetchOne(&result, "name", username); err != nil {
		log.Println(err.Error())
		return nil, false
	}

	return &result, true
}
