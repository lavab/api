package dbutils

import (
	"log"

	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
)

func GetUser(id string) (*models.User, bool) {
	var result models.User
	response, err := db.Get("users", id)
	if err == nil && !response.IsNil() {
		err := response.One(&result)
		if err != nil {
			log.Fatalln("[utils.GetUser] Error when unfolding cursor")
			return nil, false
		}
		return &result, true
	}
	return nil, false
}

func FindUserByName(username string) (*models.User, bool) {
	var result models.User
	response, err := db.GetAll("users", "name", username)
	if err == nil && response != nil && !response.IsNil() {
		err := response.One(&result)
		if err != nil {
			log.Fatalln("[utils.FindUserByName] Error when unfolding cursor")
			return nil, false
		}
		return &result, true
	}
	return nil, false
}
