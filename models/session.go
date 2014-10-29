package models

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/models/base"
)

// Session TODO
type Session struct {
	base.Expiring
	base.Resource
}

// CurrentSession returns the current request's session object
func CurrentSession(r *http.Request) *Session {
	session, ok := context.Get(r, "session").(*Session)
	if !ok {
		log.Fatalln("Session data in gorilla/context was not found or malformed.")
	}
	return session
}
