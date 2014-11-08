package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/utils"
)

var sessions = dbutils.Sessions

// AuthWrapper is an auth middleware using the "Auth" header
// The session object gets saved in the gorilla/context map, use context.Get("session") to fetch it
func AuthWrapper(next handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Auth")
		if authToken == "" {
			utils.ErrorResponse(w, 401, "Missing auth token", "")
			return
		}
		session, ok := sessions.GetSession(authToken)
		if !ok {
			utils.ErrorResponse(w, 401, "Invalid auth token", "")
			return
		}
		if session.HasExpired() {
			utils.ErrorResponse(w, 419, "Authentication token has expired", "Session has expired on "+session.ExpDate)
			sessions.DeleteId(session.ID)
			return
		}

		context.Set(r, "session", session)
		next(w, r)
	}
}
