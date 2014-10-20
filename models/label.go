package models

type Label struct {
	ID            string `json:"id" gorethink:"id"`
	UserID        string `json:"user_id" gorethink:"user_id"`
	Name          string `json:"name" gorethink:"name"`
	IsSystem      bool   `json:"is_system" gorethink:"is_system" default:"false"`
	IsVisible     bool   `json:"is_visible" gorethink:"is_visible" default:"false"`
	ThreadsUnread int    `json:"threads_unread"`
	ThreadsTotal  int    `json:"threads_total"`
	EmailsUnread  int    `json:"emails_unread" gorethink:"emails_unread"`
	EmailsTotal   int    `json:"emails_total" gorethink:"emails_total"`
}
