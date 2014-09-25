package db

import "fmt"
import "github.com/lavab/api/models"

// UsersTable TODO
var UsersTable = map[string]models.User{}

// GetUser TODO
func GetUser(username string) (models.User, error) {
	v, ok := UsersTable[username]
	if !ok {
		return models.User{}, fmt.Errorf("User not found")
	}
	return v, nil
}

// CreateUser TODO
func CreateUser(data models.User) error {

	return nil
}
