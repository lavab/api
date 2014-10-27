package routes

import (
	"net/http"

	"github.com/lavab/api/utils"
)

// Me TODO
func Me(w http.ResponseWriter, r *http.Request) {
	session := utils.CurrentSession(r)
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
}
