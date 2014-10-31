package models

import "github.com/lavab/api/models/base"

// Email is the cornerstone of our application.
// TODO mime info
type Email struct {
	base.Resource

	// AttachmentsIDs is a slice of the FileIDs associated with this email
	// For uploading attachments see `POST /upload`
	AttachmentIDs []string `json:"attachments" gorethink:"attachments"`

	// Body contains all the data needed to send this email
	Body base.Encrypted `json:"body" gorethink:"body"`

	LabelIDs []string `json:"label_ids" gorethink:"label_ids"`

	// Preview contains the encrypted preview information (needed to show a list of emails)
	// Example: Headers []string, Body string,
	// 		Headers       []string
	// 		Body          string
	// 		Snippet       string
	Preview base.Encrypted `json:"preview" gorethink:"preview"`

	// ThreadID
	ThreadID string `json:"thread_id" gorethink:"thread_id"`
}
