package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/db"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

const SessionDurationInHours = 72

// Login gets a username and password and returns a session token on success
func Login(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	account, ok := dbutils.FindAccountByUsername(username)
	if !ok || account == nil || !utils.BcryptVerify(account.Password, password) {
		utils.ErrorResponse(w, 403, "Wrong username or password",
			fmt.Sprintf("account: %+v", account))
		return
	}

	// TODO check number of sessions for the current account here
	session := models.Token{Resource: models.MakeResource(account.ID, "")}
	session.ExpireAfterNHours(SessionDurationInHours)
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

	if _, ok := dbutils.FindAccountByUsername(username); ok {
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

	account := models.Account{
		Resource: models.MakeResource(utils.UUID(), username),
		Password: string(hash),
	}

	if err := db.Insert("account", account); err != nil {
		utils.ErrorResponse(w, 500, "Internal server error",
			fmt.Sprintf("Couldn't insert %+v to database", account))
	}

	utils.JSONResponse(w, 201, map[string]interface{}{
		"message": "Signup successful",
		"success": true,
		"data":    account,
	})
}

// Logout destroys the current session token
func Logout(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*models.Token)
	if err := db.Delete("sessions", session.ID); err != nil {
		utils.ErrorResponse(w, 500, "Internal server error",
			fmt.Sprint("Couldn't delete session %v. %v", session, err))
	}
	utils.JSONResponse(w, 410, map[string]interface{}{
		"message": fmt.Sprintf("Successfully logged out", session.Owner),
		"success": true,
		"deleted": session.ID,
	})
}
