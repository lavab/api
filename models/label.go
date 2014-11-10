package models

// TODO: nested labels?

// Label is what IMAP calls folders, some providers call tags, and what we (and Gmail) call labels.
// It's both a simple way for users to organise their emails, but also a way to provide classic folder
// functionality (inbox, spam, drafts, etc).
// Examples:
//		* star an email: add the "starred" label
//		* archive an email: remove the "inbox" label
//		* delete an email: apply the "deleted" label (and cue for deletion)
type Label struct {
	Resource

	// Builtin indicates whether a label is created/needed by the system.
	// Examples: inbox, trash, spam, drafts, starred, etc.
	Builtin bool `json:"builtin" gorethink:"builtin"`

	// EmailsUnread is the number of unread emails that have a particular label applied.
	// Storing this for each label eliminates the need of db lookups for this commonly needed information.
	EmailsUnread int `json:"emails_unread" gorethink:"emails_unread"`

	// EmailsTotal is the number of emails that have a particular label applied.
	// Storing this for each label eliminates the need of db lookups for this commonly needed information.
	EmailsTotal int `json:"emails_total" gorethink:"emails_total"`
}
