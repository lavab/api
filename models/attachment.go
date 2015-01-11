package models

// Attachment is an encrypted file stored by Lavaboom
type Attachment struct {
	Encrypted
	Resource

	// Mime is the Internet media type of the attachment
	// Format: "type/subtype" – more info: en.wikipedia.org/wiki/Internet_media_type
	MIME string `json:"mime" gorethink:"mime"`

	// Size is the size of the file in bytes i.e. len(file.Data)
	Size int `json:"size" gorethink:"size"`
}
