package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/db"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/models/base"
	"github.com/lavab/api/utils"
)

type TokensGetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Created string `json:"created,omitempty"`
	Expires string `json:"expires,omitempty"`
}

func TokensGet(w http.ResponseWriter, r *http.Request) {
	// Fetch the current session from the database
	session := models.CurrentSession(r)

	// Respond with the token information
	utils.JSONResponse(200, &TokensGetResponse{
		Success: true,
		Created: session.DateCreated,
		Expires: session.ExpirationDate,
	})
}

type TokensCreateRequest struct {
	Username string `json:"username" schema:"username"`
	Password string `json:"password" schema:"password"`
}

type TokensCreateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Token   *models.Session `json:"token,omitempty"`
}

func TokensCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input TokensCreateRequest
	err := utils.ParseRequest(r, input)
	if err != nil {
		env.G.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warning("Unable to decode a request")

		utils.JSONResponse(w, 409, &TokensCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Authenticate the user
	user, ok := dbutils.FindUserByName(username)
	if !ok || user == nil || !utils.BcryptVerify(user.Password, password) {
		utils.JSONResponse(w, 403, &TokensCreateResponse{
			Success: false,
			Message: "Wrong username or password",
		})
		return
	}

	// Calculate the expiration date
	expDate := utils.HoursFromNowString(env.G.Config.SessionDuration)

	// Create a new token
	token := &models.Session{
		Expiring: base.Expiring{expDate},
		Resource: base.MakeResource(user.ID, ""),
		Name:     "Auth token expiring on " + expDate,
	}

	// Insert int into the database
	db.Insert("sessions", token)

	// Respond with the freshly created token
	utils.JSONResponse(w, 201, &TokensCreateResponse{
		Success: true,
		Message: "Authentication successful",
		Token:   token,
	})
}

// Logout destroys the current session token
type TokensDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func TokensDelete(c *web.C, w http.ResponseWriter, r *http.Request) {
	// Get the session from the middleware
	session := c.Env["session"].(*models.Session)

	// Delete it from the database
	if err := db.Delete("sessions", session.ID); err != nil {
		env.G.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to delete a session")

		utils.JSONResponse(w, 500, &TokensDeleteResponse{
			Success: true,
			Message: "Internal server error - TO/DE/01",
		})
		return
	}

	utils.JSONResponse(w, 200, &TokensDeleteResponse{
		Success: true,
		Message: "Successfully logged out",
	})
}
