package routes

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// TokensGetResponse contains the result of the TokensGet request.
type TokensGetResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	Expires *time.Time `json:"expires,omitempty"`
}

// TokensGet returns information about the current token.
func TokensGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Respond with the token information
	utils.JSONResponse(w, 200, &TokensGetResponse{
		Success: true,
		Created: &session.DateCreated,
		Expires: &session.ExpiryDate,
	})
}

// TokensCreateRequest contains the input for the TokensCreate endpoint.
type TokensCreateRequest struct {
	Username string `json:"username" schema:"username"`
	Password string `json:"password" schema:"password"`
	Type     string `json:"type" schema:"type"`
}

// TokensCreateResponse contains the result of the TokensCreate request.
type TokensCreateResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Token   *models.Token `json:"token,omitempty"`
}

// TokensCreate allows logging in to an account.
func TokensCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input TokensCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 409, &TokensCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// We can only create "auth" tokens now
	if input.Type != "auth" {
		utils.JSONResponse(w, 409, &TokensCreateResponse{
			Success: false,
			Message: "Only auth tokens are implemented",
		})
		return
	}

	// Check if account exists
	user, err := env.Accounts.FindAccountByName(input.Username)
	if err != nil {
		utils.JSONResponse(w, 403, &TokensCreateResponse{
			Success: false,
			Message: "Wrong username or password",
		})
		return
	}

	// Verify the password
	valid, updated, err := user.VerifyPassword(input.Password)
	if err != nil || !valid {
		utils.JSONResponse(w, 403, &TokensCreateResponse{
			Success: false,
			Message: "Wrong username or password",
		})
		return
	}

	// Update the user if password was updated
	if updated {
		err := env.Accounts.UpdateID(user.ID, user)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"user":  user.Name,
				"error": err,
			}).Error("Could not update user")
		}
	}

	// Calculate the expiry date
	expDate := time.Now().Add(time.Hour * time.Duration(env.Config.SessionDuration))

	// Create a new token
	token := &models.Token{
		Expiring: models.Expiring{expDate},
		Resource: models.MakeResource(user.ID, "Auth token expiring on "+expDate.Format(time.RFC3339)),
		Type:     input.Type,
	}

	// Insert the token into the cache
	err = env.TokensCache.SetToken(token)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"user":  user.Name,
			"error": err,
		}).Error("Could add token to the store")
	}

	// Insert int into the database
	env.Tokens.Insert(token)

	// Respond with the freshly created token
	utils.JSONResponse(w, 201, &TokensCreateResponse{
		Success: true,
		Message: "Authentication successful",
		Token:   token,
	})
}

// TokensDeleteResponse contains the result of the TokensDelete request.
type TokensDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TokensDelete destroys either the current auth token or the one passed as an URL param
func TokensDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Initialize
	var (
		token *models.Token
		err   error
	)

	id, ok := c.URLParams["id"]
	if !ok || id == "" {
		// Get the token from the middleware
		token = c.Env["session"].(*models.Token)
	} else {
		token, err = env.Tokens.GetToken(id)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err,
				"id":    id,
			}).Warn("Unable to find the token")

			utils.JSONResponse(w, 500, &TokensDeleteResponse{
				Success: true,
				Message: "Internal server error - TO/DE/01",
			})
			return
		}
	}

	// Delete it from the database
	if err := env.Tokens.DeleteID(token.ID); err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to delete a token")

		utils.JSONResponse(w, 500, &TokensDeleteResponse{
			Success: true,
			Message: "Internal server error - TO/DE/02",
		})
		return
	}

	//There is a lot of code repetition here ?
	if err := env.TokensCache.InvalidateToken(session.ID); err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to delete a session from cache")

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
