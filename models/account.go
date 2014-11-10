package models

// Account stores essential data for a Lavaboom user, and is thus not encrypted.
type Account struct {
	Resource

	// Billing is a struct containing billing information.
	// TODO Work in progress
	Billing BillingData `json:"billing" gorethink:"billing"`

	// Password is the password used to login to the account.
	// It's hashed and salted using a cryptographically strong method (bcrypt|scrypt).
	Password string `json:"-"  gorethink:"password"`

	// PgpExpDate is an RFC3339-encoded string containing the expiry date of the user's public key
	PgpExpDate string `json:"pgp_exp_date" gorethink:"pgp_exp_date"`

	// PgpFingerprint is a SHA-512 hash of the user's public key
	PgpFingerprint string `json:"pgp_fingerprint" gorethink:"pgp_fingerprint"`

	// PgpPublicKey is a copy of the user's current public key. It can also be found in the 'keys' db.
	PgpPublicKey string `json:"pgp_public_key" gorethink:"pgp_public_key"`

	// Settings contains data needed to customize the user experience.
	// TODO Work in progress
	Settings SettingsData `json:"settings" gorethink:"settings"`

	// Type is the account type.
	// Examples (work in progress):
	//		* beta: while in beta these are full accounts; after beta, these are normal accounts with special privileges
	//		* std: standard, free account
	//		* premium: premium account
	//		* superuser: Lavaboom staff
	Type string `json:"type" gorethink:"type"`
}

// SettingsData TODO
type SettingsData struct {
}

// BillingData TODO
type BillingData struct {
}
