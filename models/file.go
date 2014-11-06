package models

// File is an encrypted file stored by Lavaboom
type File struct {
	Encrypted
	Resource

	// Mime is the Internet media type of the file
	// Check out: http://en.wikipedia.org/wiki/Internet_media_type
	Mime string `json:"mime" gorethink:"mime"`

	// Size is the size of the file in bytes i.e. len(file.Data)
	Size int `json:"size" gorethink:"size"`

	// Type is the generic type of the file
	// Possible values: `file`, `audio`, `video`, `pdf`, `text`, `binary`
	Type string `json:"type" gorethink:"type"`
}

// Avatar is a picture used to identify contacts
type Avatar struct {
	Resource
	Large  File `json:"data" gorethink:"data"`
	Medium File `json:"thumb_small" gorethink:"thumb_small"`
	Small  File `json:"thumb_large" gorethink:"thumb_large"`
}
