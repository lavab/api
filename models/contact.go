package models

type Contact struct {
	ID     string            `json:"id" gorethink:"id"`
	Name   string            `json:"name" gorethink:"name"`
	Email  []EmailField      `json:"email" gorethink:"email"`
	Fields map[string]string `json:"fields" gorethink:"fields,omitempty"`
}

type EmailField struct {
	Address   string
	PgpKey    string
	PgpFinger string
}
