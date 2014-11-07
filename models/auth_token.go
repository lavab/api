package models

// AuthToken is a UUID used for user authentication, stored in the "auth_tokens" database
type AuthToken struct {
	Expiring
	Resource
}
