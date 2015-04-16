package models

// Thread is the data model for a list of emails, usually making up a conversation.
type Thread struct {
	Resource

	// Emails is a list of email IDs belonging to this thread
	Emails []string `json:"emails" gorethink:"emails"`

	// Labels is a list of label IDs assigned to this thread.
	// Note that emails lack this functionality. This way you can't only archive part of a thread.
	Labels []string `json:"labels" gorethink:"labels"`

	// Members is a slice containing userIDs or email addresses for all members of the thread
	Members []string `json:"members" gorethink:"members"`

	IsRead   bool   `json:"is_read" gorethink:"is_read"`
	LastRead string `json:"last_read" gorethink:"last_read"`

	Manifest string `json:"manifest,omitempty" gorethink:"manifest"`

	// SHA256 hash of the raw subject without prefixes
	SubjectHash string `json:"subject_hash" gorethink:"subject_hash"`

	// all, some, none
	Secure string `json:"secure" gorethink:"secure"`
}
