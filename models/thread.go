package models

type Thread struct {
	ID      string  `json:"id" gorethink:"id"`
	UserID  string  `json:"user_id" gorethink:"user_id"`
	Snippet string  `json:"snippet" gorethink:"snippet"`
	Changed string  `json:"changed" gorethink:"changed"`
	Emails  []Email `json:"emails" gorethink:"emails"`
}
