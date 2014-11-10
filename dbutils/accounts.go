package dbutils

import (
	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
)

func GetAccount(id string) (*models.Account, bool) {
	var result models.Account
	response, err := db.Get("accounts", id)
	if err == nil && !response.IsNil() {
		err := response.One(&result)
		if err != nil {
			return nil, false
		}
		return &result, true
	}
	return nil, false
}

func FindAccountByUsername(username string) (*models.Account, bool) {
	var result models.Account
	response, err := db.GetAll("accounts", "name", username)
	if err == nil && response != nil && !response.IsNil() {
		err := response.One(&result)
		if err != nil {
			return nil, false
		}
		return &result, true
	}
	return nil, false
}
