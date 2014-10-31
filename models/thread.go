package models

import "github.com/lavab/api/models/base"

// Thread is the data model for a conversation.
type Thread struct {
	base.Resource

	// Emails is a slice containing email IDs.
	// Ideally the array should be ordered by creation date, newest first
	EmailIDs []string `json:"email_ids" gorethink:"email_ids"`

	// Members is a slice containing userIDs or email addresses for all members of the thread
	Members []string `json:"members" gorethink:"members"`

	// Preview contains the thread details to be shown in the list of emails
	// This should be
	Preview base.Encrypted `json:"preview" gorethink:"preview"`
}
