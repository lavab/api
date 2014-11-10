package models

// Contact is the data model for a contact.
type Contact struct {
	Encrypted
	Resource

	// ProfilePicture is an encrypted picture associated with a contact.
	ProfilePicture File `json:"profile_picture" gorethink:"profile_picture"`
}
