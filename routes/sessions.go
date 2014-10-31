package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/db"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/models"
	"github.com/lavab/api/models/base"
	"github.com/lavab/api/utils"
)

const SessionDurationInHours = 72

// Login gets a username and password and returns a session token on success
func Login(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	user, ok := dbutils.FindUserByName(username)
	if !ok || user == nil || !utils.BcryptVerify(user.Password, password) {
		utils.ErrorResponse(w, 403, "Wrong username or password",
			fmt.Sprintf("user: %+v", user))
		return
	}

	// TODO check number of sessions for the current user here
	session := models.Session{
		Expiring: base.Expiring{utils.HoursFromNowString(SessionDurationInHours)},
		Resource: base.MakeResource(user.ID, ""),
	}
	session.Name = fmt.Sprintf("Auth session expiring on %s", session.ExpDate)
	db.Insert("sessions", session)

	utils.JSONResponse(w, 200, map[string]interface{}{
		"message": "Authentication successful",
		"success": true,
		"session": session,
	})
}

// Signup gets a username and password and creates a user account on success
func Signup(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	// regt := r.FormValue("reg_token")

	if _, ok := dbutils.FindUserByName(username); ok {
		utils.ErrorResponse(w, 409, "Username already exists", "")
		return
	}

	hash, err := utils.BcryptHash(password)
	if err != nil {
		msg := "Bcrypt hashing has failed"
		utils.ErrorResponse(w, 500, "Internal server error", msg)
		log.Fatalln(msg)
	}

	// TODO: sanitize user name (i.e. remove caps, periods)

	user := models.User{
		Resource: base.MakeResource(utils.UUID(), username),
		Password: string(hash),
	}

	if err := db.Insert("users", user); err != nil {
		utils.ErrorResponse(w, 500, "Internal server error",
			fmt.Sprintf("Couldn't insert %+v to database", user))
	}

	utils.JSONResponse(w, 201, map[string]interface{}{
		"message": "Signup successful",
		"success": true,
		"data":    user,
	})
}

// Logout destroys the current session token
func Logout(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*models.Session)
	if err := db.Delete("sessions", session.ID); err != nil {
		utils.ErrorResponse(w, 500, "Internal server error",
			fmt.Sprint("Couldn't delete session %v. %v", session, err))
	}
	utils.JSONResponse(w, 410, map[string]interface{}{
		"message": fmt.Sprintf("Successfully logged out", session.UserID),
		"success": true,
		"deleted": session.ID,
	})
}
