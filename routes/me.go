package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/db"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// Me returns information about the current user (more exactly, a JSONized models.User)
func Me(w http.ResponseWriter, r *http.Request) {
	session, _ := context.Get(r, "session").(*models.AuthToken)
	user, ok := dbutils.GetUser(session.AccountID)
	if !ok {
		debug := fmt.Sprintf("Session %s was deleted", session.ID)
		if err := db.Delete("sessions", session.ID); err != nil {
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
