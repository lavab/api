package models

type Contact struct {
	Resource
	Encoding  string `json:"encoding" gorethink:""`
	Format    string `json:"format" gorethink:"format"`
	PgpFinger string `json:"pgp_finger" gorethink:"pgp_finger"`
	Raw       []byte `json:"raw" gorethink:"raw"`
	Version   string `json:"" gorethink:""`
}
