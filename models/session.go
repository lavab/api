package models

// Session TODO
type Session struct {
	User      string `json:"user"`
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	ExpDate   string `json:"exp_date"`
}
