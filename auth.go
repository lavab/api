package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/utils"
)

// AuthWrapper is a middleware that checks for a session token and
// if present saves it in gorilla/context under "session".
// Session tokens are passed through the HTTP "Auth" header.
func AuthWrapper(next handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Auth")
		if authToken == "" {
			utils.ErrorResponse(w, 401, "Missing auth token", "")
			return
		}
		if session, ok := dbutils.GetSession(authToken); ok {
			context.Set(r, "session", session)
			next(w, r)
		} else {
			utils.ErrorResponse(w, 401, "Invalid auth token", "")
		}
	}
}
