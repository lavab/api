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

	UnreadThreadsCount int `json:"unread_threads_count" gorethink:"unread_threads_count"`
	TotalThreadsCount  int `json:"total_threads_count" gorethink:"total_threads_count"`
}
