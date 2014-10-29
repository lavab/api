package base

import "log"

// Encrypted is the base struct for PGP-encrypted resources
type Encrypted struct {
	// Encoding tells the reader how to decode the data; can be "json", "protobuf", maybe more in the future
	Encoding string `json:"encoding" gorethink:"encoding"`

	// PgpFinger is the fingerprint of the PGP public key used to encrypt the data. Although the obvious
	// place to look for the public key is in User.Pgp, if we store this independently we'll be able to
	// detect whether the user has replaced their key pair, and subsequently depracate this data.
	PgpFinger string `json:"pgp_finger" gorethink:"pgp_finger"`

	// Schema show how the data is structured; it's consider implicit knowledge and thus not persisted to db.
	// Lavaboom is a zero knowledge email provider, thus the less the server knows, the better.
	// It's the client's responsability to encode, decode, save and retrieve the data correctly.
	Schema map[string]interface{} `json:"schema" gorethink:"-"`

	// Raw is the raw, PGP-encrypted data
	Raw []byte `json:"raw" gorethink:"raw"`

	// Version is the schema version
	// The format is "NAME MAJOR.MINOR". If you need the actual values, use
	// Versioning the schema we can change the schema without changing
	// the data models. Minor versions can't break the API, they can only add fields (see protobuf).
	Version string `json:"version" gorethink:"version"`

	// VersionMajor is the major component of the schema version.
	// Schemas with the same major version should be compatible.
	VersionMajor int `json:"version_major" gorethink:"version_major"`

	// VersionMinor is the minor component of the schema version.
	// Schemas with different minor versions should be compatible.
	VersionMinor int `json:"version_minor" gorethink:"version_minor"`
}

// IsCompatibleWith checks whether the schema versions of two Encrypted objects are compatible
func (v *Encrypted) IsCompatibleWith(res *Encrypted) bool {
	if v == nil || res == nil {
		log.Printf("[models.IsCompatibleWith] %+v or %+v were nil\n", v, res)
		return false
	}
	if v.VersionMajor == res.VersionMajor {
		return true
	}
	return false
}
