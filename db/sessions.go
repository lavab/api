package db

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/lavab/api/models"
	"github.com/lavab/api/util"
)

// SessionsStore maps tokens to session objects
var SessionsStore = map[string]models.Session{}

// GetSession TODO
func GetSession(token string) (models.Session, error) {
	if v, ok := SessionsStore[token]; ok {
		return v, nil
	}
	return models.Session{}, fmt.Errorf("Session token not found")
}

// CreateSession TODO
func CreateSession(user, userID, userAgent string) (string, error) {
	token := uuid.New()
	SessionsStore[token] = models.Session{
		User:      user,
		UserID:    userID,
		UserAgent: userAgent,
		Expires:   util.HoursFromNow(80),
	}
	return token, nil
}
