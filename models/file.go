package models

// File is an encrypted file stored by Lavaboom
type File struct {
	Encrypted
	Resource

	// Mime is the Internet media type of the file
	// Format: "type/subtype" – more info: en.wikipedia.org/wiki/Internet_media_type
	Mime string `json:"mime" gorethink:"mime"`

	// Size is the size of the file in bytes i.e. len(file.Data)
	Size int `json:"size" gorethink:"size"`
}
