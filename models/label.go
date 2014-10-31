package models

import "github.com/lavab/api/models/base"

type Label struct {
	base.Resource
	EmailsUnread  int  `json:"emails_unread" gorethink:"emails_unread"`
	EmailsTotal   int  `json:"emails_total" gorethink:"emails_total"`
	Hidden        bool `json:"hidden" gorethink:"hidden"`
	Immutable     bool `json:"immutable" gorethink:"immutable"`
	ThreadsUnread int  `json:"threads_unread"`
	ThreadsTotal  int  `json:"threads_total"`
}
