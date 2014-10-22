package models

type Contact struct {
	ID     string            `json:"id" gorethink:"id"`
	Name   string            `json:"name" gorethink:"name"`
	Email  []EmailField      `json:"email" gorethink:"email"`
	Fields map[string]string `json:"fields" gorethink:"fields,omitempty"`
}

type EmailField struct {
	Address   string `json:"address" gorethink:"address"`
	PgpKey    string `json:"pgp_key" gorethink:"pgp_key"`
	PgpFinger string `json:"pgp_finger" gorethink:"pgp_finger"`
}
