package models

import "github.com/lavab/api/models/base"

type Thread struct {
	base.Resource
	// TODO add members?
	Snippet string  `json:"snippet" gorethink:"snippet"`
	Emails  []Email `json:"emails" gorethink:"emails"`
}
