package models

// User stores essential data for a Lavaboom user, and is thus not encrypted.
type User struct {
	Resource

	// Billing is a struct containing billing information.
	Billing BillingData `json:"billing" gorethink:"billing"`

	// Password is the actual user password, hashed using bcrypt.
	Password string `json:"-"  gorethink:"password"`

	// PgpExpDate is an RFC3339-encoded string containing the expiry date of the user's public key
	PgpExpDate string `json:"pgp_exp_date" gorethink:"pgp_exp_date"`

	// PgpFingerprint is a SHA-512 hash of the user's public key
	PgpFingerprint string `json:"pgp_fingerprint" gorethink:"pgp_fingerprint"`

	// PgpPublicKey is a copy of the user's current public key. It can also be found in the 'keys' db.
	PgpPublicKey string `json:"pgp_public_key" gorethink:"pgp_public_key"`

	// Settings is a struct containing app configuration data.
	Settings SettingsData `json:"settings" gorethink:"settings"`

	// Type is the user type (free, beta, premium, etc)
	Type string `json:"type" gorethink:"type"`
}

// SettingsData TODO
type SettingsData struct {
}

// BillingData TODO
type BillingData struct {
}
