package models

import (
	"time"

	"github.com/dchest/uniuri"
)

// Resource is the base type for API resources.
type Resource struct {
	// ID is the resources ID, used as a primary key by the db.
	// For some resources (invites, auth tokens) this is also the data itself.
	ID string `json:"id" gorethink:"id"`

	// DateCreated is, shockingly, the time when the resource was created.
	DateCreated time.Time `json:"date_created" gorethink:"date_created"`

	// DateModified records the time of the last change of the resource.
	DateModified time.Time `json:"date_modified" gorethink:"date_modified"`

	// Name is the human-friendly name of the resource. It can either be essential (e.g. Account.Name) or optional.
	Name string `json:"name" gorethink:"name,omitempty"`

	// Owner is the ID of the account that owns this resource.
	Owner string `json:"owner" gorethink:"owner"`
}

// MakeResource creates a new Resource object with sane defaults.
func MakeResource(ownerID, name string) Resource {
	t := time.Now()
	return Resource{
		ID:           uniuri.NewLen(uniuri.UUIDLen),
		DateModified: t,
		DateCreated:  t,
		Name:         name,
		Owner:        ownerID,
	}
}

// Touch sets the time the resource was last modified to time.Now().
// For convenience (e.g. chaining) it also returns the resource pointer.
func (r *Resource) Touch() *Resource {
	r.DateModified = time.Now()
	return r
}
