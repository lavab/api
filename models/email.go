package models

// Email is the cornerstone of our application.
// TODO mime info
type Email struct {
	Resource

	// Kind of the email. Value is either sent or received.
	Kind string `json:"kind" gorethink:"kind"`

	From []string `json:"from" gorethink:"from"`

	// Who is supposed to receive the email / what email received it.
	To []string `json:"to" gorethink:"to"`

	// AttachmentsIDs is a slice of the FileIDs associated with this email
	// For uploading attachments see `POST /upload`
	AttachmentIDs []string `json:"attachments" gorethink:"attachments"`

	// Body contains all the data needed to send this email
	Body Encrypted `json:"body" gorethink:"body"`

	LabelIDs []string `json:"label_ids" gorethink:"label_ids"`

	// Preview contains the encrypted preview information (needed to show a list of emails)
	// Example: Headers []string, Body string,
	// 		Headers       []string
	// 		Body          string
	// 		Snippet       string
	Preview Encrypted `json:"preview" gorethink:"preview"`

	// ThreadID
	ThreadID string `json:"thread_id" gorethink:"thread_id"`

	Status string `json:"status" gorethink:"status"`

	IsRead string `json:"is_read" gorethink:"is_read"`
}
