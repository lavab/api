package models

// File is an encrypted file stored by Lavaboom
type File struct {
	Resource

	Meta interface{} `json:"meta" gorethink:"meta"`
	Body []byte      `json:"body" gorethink:"body"`
	Tags []string    `json:"tags" gorethink:"tags"`
}
