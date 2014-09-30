package db

import (
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// TODO replace with rethinkdb
var sessionStore = map[string]models.Session{}

// Session TODO
func Session(token string) (models.Session, error) {
	v, ok := sessionStore[token]
	if !ok {
		return models.Session{}, fmt.Errorf("Session token not found")
	}
	expDate, err := time.Parse(time.RFC3339, v.ExpDate)
	if err == nil && time.Now().UTC().After(expDate) {
		DeleteSession(token)
		return models.Session{}, fmt.Errorf("Session expired")
	}
	return v, nil
}

// CreateSession TODO
func CreateSession(user string, hours int) (string, error) {
	token := uuid.New()
	sessionStore[token] = models.Session{
		User:    user,
		ExpDate: utils.HoursFromNow(hours),
	}
	return token, nil
}

// DeleteSession TODO
func DeleteSession(token string) error {
	delete(sessionStore, token)
	return nil
}
