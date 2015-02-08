package models

// File is an encrypted file stored by Lavaboom
type File struct {
	Encrypted
	Resource

	Metadata Encrypted `json:"encrypted" gorethink:"encrypted"`
	Body     Encrypted `json:"body" gorethink:"body"`
}
