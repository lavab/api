package models

// Thread is the data model for a conversation.
type Thread struct {
	Resource

	// Emails is an array of email IDs belonging to this thread
	Emails []string `json:"emails" gorethink:"emails"`

	// Members is a slice containing userIDs or email addresses for all members of the thread
	Members []string `json:"members" gorethink:"members"`
}
