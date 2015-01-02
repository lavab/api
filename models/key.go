package models

type Key struct {
	Resource // ID is the fingerprint, Name is empty
	Expiring // ExpiryDate is either empty or expiring, user can set it

	//Body []byte `json:"body" gorethink:"body"` // Raw key contents

	Headers     map[string]string `json:"headers" gorethink:"headers"`           // Headers passed with the key
	Algorithm   string            `json:"algorithm" gorethink:"algorithm"`       // Algorithm of the key
	Length      uint16            `json:"length" gorethink:"length"`             // Length of the key
	Key         string            `json:"key" gorethink:"key"`                   // Armor-encoded key
	KeyID       string            `json:"key_id" gorethink:"key_id"`             // PGP key ID
	KeyIDShort  string            `json:"key_id_short" gorethink:"key_id_short"` // Shorter version of above
	Reliability int               `json:"reliability" gorethink:"reliability"`   // Reliability algorithm cached result
	MasterKey   string            `json:"master_key" gorethink:"mater_key"`      // MasterKey's ID - no idea how it works
}
