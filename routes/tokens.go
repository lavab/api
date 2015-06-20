package routes

import (
	"net/http"
	"time"

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
			utils.JSONResponse(w, 404, utils.NewError(
				utils.TokensGetUnableToGet, err, false,
			))
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.TokensCreateInvalidInput, err, false,
		))
		return
	}

	// We can only create "auth" tokens now
	if input.Type != "auth" {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.TokensCreateInvalidType, "Only auth tokens are implemented", false,
		))
		return
	}

	input.Username = utils.RemoveDots(
		utils.NormalizeUsername(input.Username),
	)

	// Check if account exists
	user, err := env.Accounts.FindAccountByName(input.Username)
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.TokensCreateUnableToGetAccount, err, false,
		))
		return
	}

	// "registered" accounts can't log in
	if user.Status == "registered" {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.TokensCreateInvalidStatus, "Your account is not confirmed", false,
		))
		return
	}

	// Verify the password
	valid, updated, err := user.VerifyPassword(input.Password)
	if err != nil || !valid {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.TokensCreateInvalidPassword, err, false,
		))
		return
	}

	// Update the user if password was updated
	if updated {
		user.DateModified = time.Now()
		err := env.Accounts.UpdateID(user.ID, user)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.TokensCreateUnableToUpdate, err, true,
			))
			return
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
	if err := env.Tokens.Insert(token); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.TokensCreateUnableToInsert, err, true,
		))
		return
	}

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
			utils.JSONResponse(w, 404, utils.NewError(
				utils.TokensDeleteUnableToGet, err, false,
			))
			return
		}
	}

	// Delete it from the database
	if err := env.Tokens.DeleteID(token.ID); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.TokensDeleteUnableToDelete, err, true,
		))
		return
	}

	utils.JSONResponse(w, 200, &TokensDeleteResponse{
		Success: true,
		Message: "Successfully logged out",
	})
}
