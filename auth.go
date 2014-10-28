package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/db"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/utils"
)

// AuthWrapper is an auth middleware using the "Auth" header
// The session object gets saved in the gorilla/context map, use context.Get("session") to fetch it
func AuthWrapper(next handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Auth")
		if authToken == "" {
			utils.ErrorResponse(w, 401, "Missing auth token", "")
			return
		}
		session, ok := dbutils.GetSession(authToken)
		if !ok {
			utils.ErrorResponse(w, 401, "Invalid auth token", "")
			return
		}
		if session.HasExpired() {
			utils.ErrorResponse(w, 419, "Authentication token has expired", "Session has expired on "+session.ExpDate)
			db.Delete("sessions", session.ID)
			return
		}

		context.Set(r, "session", session)
		next(w, r)
	}
}
