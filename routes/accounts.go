package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// AccountsListResponse contains the result of the AccountsList request.
type AccountsListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsList returns a list of accounts visible to an user
func AccountsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &AccountsListResponse{
		Success: false,
		Message: "Method not implemented",
	})
}

// AccountsCreateRequest contains the input for the AccountsCreate endpoint.
type AccountsCreateRequest struct {
	Username string `json:"username" schema:"username"`
	Password string `json:"password" schema:"password"`
}

// AccountsCreateResponse contains the output of the AccountsCreate request.
type AccountsCreateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Account *models.Account `json:"account,omitempty"`
}

// AccountsCreate creates a new account in the system.
func AccountsCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input AccountsCreateRequest
	err := utils.ParseRequest(r, input)
	if err != nil {
		env.G.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 409, &AccountsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Ensure that the user with requested username doesn't exist
	if _, err := env.G.R.Accounts.FindAccountByName(input.Username); err != nil {
		utils.JSONResponse(w, 409, &AccountsCreateResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// Try to hash the password
	hash, err := utils.BcryptHash(input.Password)
	if err != nil {
		utils.JSONResponse(w, 500, &AccountsCreateResponse{
			Success: false,
			Message: "Internal server error - AC/CR/01",
		})

		env.G.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to hash a password")
		return
	}

	// TODO: sanitize user name (i.e. remove caps, periods)

	// Create a new user object
	account := &models.Account{
		Resource: models.MakeResource("", input.Username),
		Password: string(hash),
	}

	// Try to save it in the database
	if err := env.G.R.Accounts.Insert(account); err != nil {
		utils.JSONResponse(w, 500, &AccountsCreateResponse{
			Success: false,
			Message: "Internal server error - AC/CR/02",
		})

		env.G.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not insert an user to the database")
		return
	}

	utils.JSONResponse(w, 201, &AccountsCreateResponse{
		Success: true,
		Message: "A new account was successfully created",
		Account: account,
	})
}

// AccountsGetResponse contains the result of the AccountsGet request.
type AccountsGetResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	User    *models.Account `json:"user,omitempty"`
}

// AccountsGet returns the information about the specified account
func AccountsGet(c *web.C, w http.ResponseWriter, r *http.Request) {
	// Get the account ID from the request
	id, ok := c.URLParams["id"]
	if !ok {
		utils.JSONResponse(w, 409, &AccountsGetResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsGetResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.G.R.Accounts.GetAccount(session.Owner)
	if err != nil {
		// The session refers to a non-existing user
		env.G.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err,
		}).Warn("Valid session referred to a removed account")

		// Try to remove the orphaned session
		if err := env.G.R.Tokens.DeleteID(session.ID); err != nil {
			env.G.Log.WithFields(logrus.Fields{
				"id":    session.ID,
				"error": err,
			}).Error("Unable to remove an orphaned session")
		} else {
			env.G.Log.WithFields(logrus.Fields{
				"id": session.ID,
			}).Info("Removed an orphaned session")
		}

		utils.JSONResponse(w, 410, &AccountsGetResponse{
			Success: false,
			Message: "Account disabled",
		})
		return
	}

	// Return the user struct
	utils.JSONResponse(w, 200, &AccountsGetResponse{
		Success: true,
		User:    user,
	})
}

// AccountsUpdateResponse contains the result of the AccountsUpdate request.
type AccountsUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsUpdate allows changing the account's information (password etc.)
func AccountsUpdate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &AccountsUpdateResponse{
		Success: false,
		Message: `Sorry, not implemented yet`,
	})
}

// AccountsDeleteResponse contains the result of the AccountsDelete request.
type AccountsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsDelete allows deleting an account.
func AccountsDelete(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &AccountsDeleteResponse{
		Success: false,
		Message: `Sorry, not implemented yet`,
	})
}

// AccountsWipeDataResponse contains the result of the AccountsWipeData request.
type AccountsWipeDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsWipeData allows getting rid of the all data related to the account.
func AccountsWipeData(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &AccountsWipeDataResponse{
		Success: false,
		Message: `Sorry, not implemented yet`,
	})
}

// AccountsSessionsListResponse contains the result of the AccountsSessionsList request.
type AccountsSessionsListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsSessionsList returns a list of all opened sessions.
func AccountsSessionsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &AccountsSessionsListResponse{
		Success: false,
		Message: `Sorry, not implemented yet`,
	})
}
