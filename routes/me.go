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
		"data": map[string]interface{}{
			"user":    session.User,
			"user_id": session.UserID,
		},
	})
}

// UpdateMe TODO
func UpdateMe(w http.ResponseWriter, r *http.Request) {
}

// Settings TODO
func Settings(w http.ResponseWriter, r *http.Request) {
}

// UpdateSettings TODO
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
}
