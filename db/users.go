package db

import "fmt"
import "github.com/lavab/api/models"

// UsersTable TODO
var usersTable = map[string]models.User{}

// User TODO
func User(username string) (models.User, error) {
	v, ok := usersTable[username]
	if !ok {
		return models.User{}, fmt.Errorf("User not found")
	}
	return v, nil
}

// CreateUser TODO
func CreateUser(user models.User) error {
	if _, ok := usersTable[user.Name]; !ok {
		usersTable[user.Name] = user
	} else {
		return fmt.Errorf("Username already exists")
	}
	return nil
}
