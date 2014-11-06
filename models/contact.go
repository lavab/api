package models

// Contact is the data model for a contact.
type Contact struct {
	Encrypted
	Resource

	// Picture is a profile picture
	Picture Avatar `json:"picture" gorethink:"picture"`
}
