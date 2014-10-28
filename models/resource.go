package models

import (
	"log"
	"time"

	"github.com/lavab/api/utils"
)

type Resource struct {
	ID          string `json:"id" gorethink:"id"`
	LastChanged string `json:"last_changed" gorethink:"last_changed"`
	Name        string `json:"name" gorethink:"name,omitempty"`
	UserID      string `json:"user_id" gorethink:"user_id"`
}

type Expiry struct {
	ExpDate string `json:"exp_date" gorethink:"exp_date"`
}

func MakeResource(name, userID string) Resource {
	return Resource{
		ID:          utils.UUID(),
		LastChanged: utils.TimeNowString(),
		Name:        name,
		UserID:      userID,
	}
}

func (e *Expiry) HasExpired() bool {
	t, err := time.Parse(time.RFC3339, e.ExpDate)
	if err != nil {
		log.Println("Bad format! The expiry date not RFC3339.", err)
		return true
	}
	if time.Now().UTC().After(t) {
		return true
	}
	return false
}
