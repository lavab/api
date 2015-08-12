package models

// File is an encrypted file stored by Lavaboom
type File struct {
	Encrypted
	Resource

	Meta interface{} `json:"meta" gorethink:"meta,omitempty"`
	Body []byte      `json:"body" gorethink:"body,omitempty"`
	Tags []string    `json:"tags" gorethink:"tags,omitempty"`
}
