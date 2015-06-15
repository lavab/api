package models

// File is an encrypted file stored by Lavaboom
type File struct {
	Resource

	Meta FileMeta `json:"meta" gorethink:"meta"`
	Body []byte   `json:"body" gorethink:"body"`
	Tags []string `json:"tags" gorethink:"tags"`
}

type FileMeta map[string]interface{}

func (f FileMeta) ContentType() string {
	if a, ok := f["content_type"]; ok {
		if b, ok := a.(string); ok {
			return b
		}
	}

	return ""
}
