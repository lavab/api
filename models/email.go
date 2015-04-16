package models

// Email is a message in a thread
type Email struct {
	Resource

	MessageID string `json:"message_id" gorethink:"message_id"`

	// Kind is the type of encryption used in the email:
	//  - raw      - when sending raw emails before they get sent
	//  - manifest - Manifest field is not empty,
	//  - pgpmime  - PGP/MIME format, aka everything is in body
	Kind string `json:"kind" gorethink:"kind"`

	// Unencrypted metadata information, available in both received in sent emails
	From string   `json:"from" gorethink:"from"`
	To   []string `json:"to" gorethink:"to"`
	CC   []string `json:"cc" gorethink:"cc"`

	// BCC is only visible in sent emails
	BCC []string `json:"bcc" gorethink:"bcc"`

	// Fingerprints used for body and manifest
	PGPFingerprints []string `json:"pgp_fingerprints" gorethink:"pgp_fingerprints"`

	// Files contains IDs of other files
	Files []string `json:"files" gorethink:"files"`

	// Manifest is only available in emails that were encrypted using PGP manifests
	Manifest string `json:"manifest" gorethink:"manifest"`

	// Body contains all the data needed to send this email
	Body string `json:"body" gorethink:"body"`

	// ContentType of the body in unencrypted emails
	ContentType string `json:"content_type" gorethink:"content_type"`
	ReplyTo     string `json:"reply_to" gorethink:"reply_to"`

	// Contains ID of the thread
	Thread string `json:"thread" gorethink:"thread"`

	// received or (queued|processed)
	Status string `json:"status" gorethink:"status"`
}
