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
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Token   *models.Token `json:"token,omitempty"`
}

// TokensGet returns information about the current token.
func TokensGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Initialize
	var (
		token *models.Token
		err   error
	)

	id, ok := c.URLParams["id"]
	if !ok || id == "" {
		// Get the token from the middleware
		token = c.Env["token"].(*models.Token)
	} else {
		token, err = env.Tokens.GetToken(id)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    id,
			}).Warn("Unable to find the token")

			utils.JSONResponse(w, 404, &TokensGetResponse{
				Success: false,
				Message: "Invalid token ID",
			})
			return
		}
	}

	// Respond with the token information
	utils.JSONResponse(w, 200, &TokensGetResponse{
		Success: true,
		Token:   token,
	})
}

// TokensCreateRequest contains the input for the TokensCreate endpoint.
type TokensCreateRequest struct {
	Username string `json:"username" schema:"username"`
	Password string `json:"password" schema:"password"`
	Type     string `json:"type" schema:"type"`
	Token    string `json:"token" schema:"token"`
}

// TokensCreateResponse contains the result of the TokensCreate request.
type TokensCreateResponse struct {
	Success         bool          `json:"success"`
	Message         string        `json:"message,omitempty"`
	Token           *models.Token `json:"token,omitempty"`
	FactorType      string        `json:"factor_type,omitempty"`
	FactorChallenge string        `json:"factor_challenge,omitempty"`
}

// TokensCreate allows logging in to an account.
func TokensCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input TokensCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
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

	input.Username = utils.RemoveDots(
		utils.NormalizeUsername(input.Username),
	)

	// Check if account exists
	user, err := env.Accounts.FindAccountByName(input.Username)
	if err != nil {
		utils.JSONResponse(w, 403, &TokensCreateResponse{
			Success: false,
			Message: "Wrong username or password",
		})
		return
	}

	// "registered" accounts can't log in
	if user.Status == "registered" {
		utils.JSONResponse(w, 403, &TokensCreateResponse{
			Success: false,
			Message: "Your account is not confirmed",
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
		user.DateModified = time.Now()
		err := env.Accounts.UpdateID(user.ID, user)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"user":  user.Name,
				"error": err.Error(),
			}).Error("Could not update user")

			// DO NOT RETURN!
		}
	}

	// Check for 2nd factor
	if user.FactorType != "" {
		factor, ok := env.Factors[user.FactorType]
		if ok {
			// Verify the 2FA
			verified, challenge, err := user.Verify2FA(factor, input.Token)
			if err != nil {
				utils.JSONResponse(w, 500, &TokensCreateResponse{
					Success: false,
					Message: "Internal 2FA error",
				})

				env.Log.WithFields(logrus.Fields{
					"err":    err.Error(),
					"factor": user.FactorType,
				}).Warn("2FA authentication error")
				return
			}

			// Token was probably empty. Return the challenge.
			if !verified && challenge != "" {
				utils.JSONResponse(w, 403, &TokensCreateResponse{
					Success:         false,
					Message:         "2FA token was not passed",
					FactorType:      user.FactorType,
					FactorChallenge: challenge,
				})
				return
			}

			// Token was incorrect
			if !verified {
				utils.JSONResponse(w, 403, &TokensCreateResponse{
					Success:    false,
					Message:    "Invalid token passed",
					FactorType: user.FactorType,
				})
				return
			}
		}
	}

	// Calculate the expiry date
	expDate := time.Now().Add(time.Hour * time.Duration(env.Config.SessionDuration))

	// Create a new token
	token := &models.Token{
		Expiring: models.Expiring{ExpiryDate: expDate},
		Resource: models.MakeResource(user.ID, "Auth token expiring on "+expDate.Format(time.RFC3339)),
		Type:     input.Type,
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
		token = c.Env["token"].(*models.Token)
	} else {
		token, err = env.Tokens.GetToken(id)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    id,
			}).Warn("Unable to find the token")

			utils.JSONResponse(w, 404, &TokensDeleteResponse{
				Success: false,
				Message: "Invalid token ID",
			})
			return
		}
	}

	// Delete it from the database
	if err := env.Tokens.DeleteID(token.ID); err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to delete a token")

		utils.JSONResponse(w, 500, &TokensDeleteResponse{
			Success: false,
			Message: "Internal server error - TO/DE/02",
		})
		return
	}

	utils.JSONResponse(w, 200, &TokensDeleteResponse{
		Success: true,
		Message: "Successfully logged out",
	})
}
