package models

// Token is a volatile, unique object. It can be used for user authentication, confirmations, invites, etc.
type Token struct {
	Expiring
	Resource

	// Type describes the token's purpose: auth, invite, confirm, upgrade.
	Type string `json:"type" gorethink:"type"`
}

// MakeToken creates a generic token.
func MakeToken(accountID, _type string, nHours int) Token {
	out := Token{
		Resource: MakeResource(accountID, ""),
		Type:     _type,
	}
	out.ExpireAfterNHours(nHours)
	return out
}

// Invalidate invalidates a token by adding a period (".") at the beginning of its type.
// It also shortens its expiration time.
func (t *Token) Invalidate() {
	t.Type = "." + t.Type
	t.ExpireSoon()
}

// MakeAuthToken creates an authentication token, valid for a limited time.
func MakeAuthToken(accountID string) Token {
	return MakeToken(accountID, "auth", 80)
}

// MakeInviteToken creates an invitation to create an account.
func MakeInviteToken(accountID string) Token {
	return MakeToken(accountID, "invite", 240)
}
