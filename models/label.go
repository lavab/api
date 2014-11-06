package models

// TODO: nested labels?

// Label is what IMAP calls folders, some providers call tags, and what we (and Gmail) call labels.
// It's both a simple way for users to organise their emails, but also a way to provide classic folder
// functionality (inbox, spam, drafts, etc). For example, to "archive" an email means to remove the "inbox" label.
type Label struct {
	Resource
	EmailsUnread  int  `json:"emails_unread" gorethink:"emails_unread"`
	EmailsTotal   int  `json:"emails_total" gorethink:"emails_total"`
	Hidden        bool `json:"hidden" gorethink:"hidden"`
	Immutable     bool `json:"immutable" gorethink:"immutable"`
	ThreadsUnread int  `json:"threads_unread"`
	ThreadsTotal  int  `json:"threads_total"`
}
