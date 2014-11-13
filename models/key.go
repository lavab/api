package models

type Key struct {
	Resource
	Expiring

	// ID is the fingerprint

	Image []byte `json:"image" gorethink:"image"`

	// the actual key
	Key string `json:"key" gorethink:"key"`

	OwnerName string `json:"owner_name" gorethink:"owner_name"`

	// the actual id
	KeyID      string `json:"key_id" gorethink:"key_id"`
	KeyIDShort string `json:"key_id_short" gorethink:"key_id_short"`
}
