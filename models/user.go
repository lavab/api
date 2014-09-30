package models

import "time"

// User TODO
type User struct {
	ID       string
	Name     string
	Password string
	Key      PGP
}

// PGP TODO
type PGP struct {
	Key     string
	Finger  string
	Expires time.Time
}
