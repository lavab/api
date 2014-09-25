package models

import "time"

// User TODO
type User struct {
	ID       string
	Name     string
	Password string
	Salt     string
	Key      PGP
	JWT      KeyPair
}

// PGP TODO
type PGP struct {
	Key     string
	Finger  string
	Expires time.Time
}

// KeyPair TODO
type KeyPair struct {
	Public  string
	Private string
}
