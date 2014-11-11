package models

type Key struct {
	Resource
	Expiring

	// ID is the fingerprint

	Image []byte `json:"image" gorethink:"image"`

	// the actual key
	Key string `json:"key" gorethink:"key"`

	// the actual id
	ShortID string `json:"short_id" gorethink:"short_id"`
}
