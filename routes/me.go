package routes

import (
	"encoding/json"
	"fmt"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
	"log"
	"net/http"
)

// Me returns information about the current user (more exactly, a JSONized models.User)
func Me(w http.ResponseWriter, r *http.Request) {
	session := models.CurrentSession(r)
	user, ok := users.GetUser(session.UserID)
	if !ok {
		debug := fmt.Sprintf("Session %s was deleted", session.ID)
		if err := sessions.DeleteId(session.ID); err != nil {
			debug = "Error when trying to delete session associated with inactive account"
			log.Println("[routes.Me]", debug, err)
		}
		utils.ErrorResponse(w, 410, "Account deactivated", debug)
		return
	}
	str, err := json.Marshal(user)
	if err != nil {
		debug := fmt.Sprint("Failed to marshal models.User:", user)
		log.Println("[routes.Me]", debug)
		utils.ErrorResponse(w, 500, "Internal server error", debug)
		return
	}
	fmt.Fprint(w, string(str))
}

// UpdateMe TODO
func UpdateMe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// Sessions lists all active sessions for current user
func Sessions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
