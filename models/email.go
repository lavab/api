package models

type Email struct {
	ID           string       `json:"id" gorethink:"id"`
	UserID       string       `json:"user_id" gorethink:"user_id"`
	ThreadID     string       `json:"thread_id" gorethink:"thread_id"`
	LabelIDs     []string     `json:"label_ids" gorethink:"label_ids"`
	Headers      []string     `json:"headers" gorethink:"headers"`
	Body         string       `json:"body" gorethink:"body"`
	Snippet      string       `json:"snippet" gorethink:"snippet"`
	Date         string       `json:"date" gorethink:"date"`
	SizeEstimate int          `json:"size_estimate" gorethink:"size_estimate" unit:"byte"`
	Attachments  []Attachment `json:"attachments" gorethink:"attachments"`
	PgpKeys      []string     `json:"pgp_keys" gorethink:"pgp_keys"`
	Raw          string       `json:"raw" gorethink:"raw"`
}

type Attachment struct {
	FileID       string `json:"file_id" gorethink:"file_id,omitempty"`
	Data         []byte `json:"data" gorethink:"data,omitempty"`
	MimeType     string `json:"mime_type" gorethink:"mime_type,omitempty"`
	SizeEstimate int    `json:"size" gorethink:"size" unit:"byte,omitempty"`
}
