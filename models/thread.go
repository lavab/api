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

	// Subject is the subject of the thread.
	Subject string `json:"subject" gorethink:"subject"`

	IsRead bool `json:"is_read" gorethink:"is_read"`
}
