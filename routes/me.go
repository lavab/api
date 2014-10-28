package routes

import (
	"fmt"
	"net/http"

	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// Me TODO return models.User
func Me(w http.ResponseWriter, r *http.Request) {
	session := models.CurrentSession(r)
	utils.JSONResponse(w, map[string]interface{}{
		"status": 200,
		"user": map[string]interface{}{
			"id":   session.UserID,
			"name": session.User,
		},
	})
}

// UpdateMe TODO
func UpdateMe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// Sessions lists all active sessions for current user
func Sessions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
