package models

import (
	"time"

	"github.com/lavab/api/utils"
)

// Resource is the base type for API resources.
type Resource struct {
	// ID is the resources ID, used as a primary key by the db.
	// For some resources (invites, auth tokens) this is also the data itself.
	ID string `json:"id" gorethink:"id"`

	// DateCreated is, shockingly, the date when the resource was created.
	DateCreated time.Time `json:"date_created" gorethink:"date_created"`

	// DateModified records the time of the last change of the resource.
	DateModified time.Time `json:"date_modified" gorethink:"date_modified"`

	// Name is a human-friendly description of the resource.
	// Sometimes it can be essential to the resource, e.g. the `Account.Name` field.
	Name string `json:"name" gorethink:"name,omitempty"`

	// AccountID is the ID of the user account that owns this resource.
	AccountID string `json:"user_id" gorethink:"user_id"`
}

// MakeResource creates a new Resource object with sane defaults.
func MakeResource(userID, name string) Resource {
	t := time.Now()
	return Resource{
		ID:           utils.UUID(),
		DateModified: t,
		DateCreated:  t,
		Name:         name,
		AccountID:    userID,
	}
}

// Touch sets the time the resource was last modified to time.Now().
// For convenience (e.g. chaining) it also returns the resource pointer.
func (r *Resource) Touch() *Resource {
	r.DateModified = time.Now()
	return r
}
