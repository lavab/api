package models

import (
	"time"
)

// Expiring is a base struct for resources that expires e.g. sessions.
type Expiring struct {
	// ExpiryDate indicates when an object will expire
	ExpiryDate time.Time `json:"expiry_date" gorethink:"expiry_date"`
}

// Expired checks whether an object has expired. It returns true if ExpiryDate is in the past.
func (e *Expiring) Expired() bool {
	if time.Now().UTC().After(e.ExpiryDate) {
		return true
	}
	return false
}

// ExpireAfterNHours sets the expiry date to time.Now().UTC() + n hours
func (e *Expiring) ExpireAfterNHours(n int) {
	e.ExpiryDate = time.Now().UTC().Add(time.Duration(n) * time.Hour)
}

// ExpireSoon sets the expiry date to something in the near future.
func (e *Expiring) ExpireSoon() {
	e.ExpiryDate = time.Now().UTC().Add(time.Duration(2) * time.Minute)
}
