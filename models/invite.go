package models

// Invite is a token (Invite.ID) that allows a user
type Invite struct {
	Expiring
	Resource

	// AccountCreated is the ID of the account that was created using this invite.
	AccountCreated string

	// Username is the desired username. It can be blank.
	Username string
}

// Used returns whether this invitation has been used
func (i *Invite) Used() bool {
	return i.DateCreated != i.DateModified
}
