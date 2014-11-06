package models

// Encrypted is the base struct for PGP-encrypted resources.
type Encrypted struct {
	// Encoding tells the reader how to decode the data; can be "json", "protobuf", maybe more in the future
	Encoding string `json:"encoding" gorethink:"encoding"`

	// PgpFingerprints contains the fingerprints of the PGP public keys used to encrypt the data.
	PgpFingerprints []string `json:"pgp_fingerprints" gorethink:"pgp_fingerprints"`

	// Data is the raw, PGP-encrypted data
	Data []byte `json:"raw" gorethink:"raw"`

	// Schema is the name of the schema used to encode the data
	// Examples: string, contact, email
	Schema string `json:"schema" gorethink:"schema"`

	// VersionMajor is the major component of the schema version.
	// Schemas with the same major version should be compatible.
	VersionMajor int `json:"version_major" gorethink:"version_major"`

	// VersionMinor is the minor component of the schema version.
	// Schemas with different minor versions should be compatible.
	VersionMinor int `json:"version_minor" gorethink:"version_minor"`
}
