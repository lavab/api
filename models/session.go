package models

// Session TODO
type Session struct {
	ID        string `json:"id" gorethink:"id"`
	User      string `json:"user" gorethink:"user"`
	UserID    string `json:"user_id" gorethink:"user_id"`
	UserAgent string `json:"-" gorethink:"user_agent"`
	ExpDate   string `json:"exp_date" gorethink:"exp_date"`
}
