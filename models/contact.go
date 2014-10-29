package models

import "github.com/lavab/api/models/base"

// Contact is the data model for a contact.
type Contact struct {
	base.Encrypted
	base.Resource
}
