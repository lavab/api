package models

import "github.com/lavab/api/models/base"

// Email is the cornerstone of our application.
type Email struct {
	base.Encrypted
	base.Resource
	ThreadID    string       `json:"thread_id" gorethink:"thread_id"`
	LabelIDs    []string     `json:"label_ids" gorethink:"label_ids"`
	Headers     []string     `json:"headers" gorethink:"headers"`
	Body        string       `json:"body" gorethink:"body"`
	Snippet     string       `json:"snippet" gorethink:"snippet"`
	Attachments []Attachment `json:"attachments" gorethink:"attachments"`
	PgpKeys     []string     `json:"pgp_keys" gorethink:"pgp_keys"`
}

type Attachment struct {
	FileID       string `json:"file_id" gorethink:"file_id,omitempty"`
	Data         []byte `json:"data" gorethink:"data,omitempty"`
	MimeType     string `json:"mime_type" gorethink:"mime_type,omitempty"`
	SizeEstimate int    `json:"size" gorethink:"size" unit:"byte,omitempty"`
}
