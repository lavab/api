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
	Token    string `json:"token" schema:"token"`
	Username string `json:"username" schema:"username"`
	Password string `json:"password" schema:"password"`
	AltEmail string `json:"alt_email" schema:"alt_email"`
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
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &AccountsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Detect the request type
	// 1) username + token + password     - invite
	// 2) username + password + alt_email - register with confirmation
	// 3) alt_email only                  - register for beta (add to queue)
	requestType := "unknown"
	if input.AltEmail == "" && input.Username != "" && input.Password != "" && input.Token != "" {
		requestType = "invited"
	} else if input.AltEmail != "" && input.Username != "" && input.Password != "" && input.Token != "" {
		requestType = "classic"
	} else if input.AltEmail != "" && input.Username == "" && input.Password == "" && input.Token == "" {
		requestType = "queue"
	}

	// "unknown" requests are empty and invalid
	if requestType == "invalid" {
		utils.JSONResponse(w, 400, &AccountsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Adding to queue will be implemented soon
	if requestType == "queue" {
		// Implementation awaits https://trello.com/c/SLM0qK1O/91-account-registration-queue
		utils.JSONResponse(w, 501, &AccountsCreateResponse{
			Success: false,
			Message: "Sorry, not implemented yet",
		})
		return
	}

	// Check if classic registration is enabled
	if requestType == "classic" && !env.Config.ClassicRegistration {
		utils.JSONResponse(w, 403, &AccountsCreateResponse{
			Success: false,
			Message: "Classic registration is disabled",
		})
		return
	}

	// Check "invited" for token validity
	if requestType == "invited" {
		// Fetch the token from the database
		token, err := env.Tokens.GetToken(input.Token)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Unable to fetch a registration token from the database")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation token",
			})
			return
		}

		// Ensure that the token's type is valid
		if token.Type != "invite" {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation token",
			})
			return
		}

		// Check if it's expired
		if token.Expired() {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Expired invitation token",
			})
			return
		}
	}

	// TODO: sanitize user name (i.e. remove caps, periods)

	// Both invited and classic require an unique username, so ensure that the user with requested username isn't already used
	if _, err := env.Accounts.FindAccountByName(input.Username); err == nil {
		utils.JSONResponse(w, 409, &AccountsCreateResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// Both username and password are filled, so we can create a new account.
	account := &models.Account{
		Resource: models.MakeResource("", input.Username),
		Type:     "beta",
	}

	// Set the password
	err = account.SetPassword(input.Password)
	if err != nil {
		utils.JSONResponse(w, 500, &AccountsCreateResponse{
			Success: false,
			Message: "Internal server error - AC/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to hash the password")
		return
	}

	// User won't be able to log in until the account gets verified
	if requestType == "classic" {
		account.Status = "unverified"
	}

	// Set the status to invited, because of stats
	if requestType == "invited" {
		account.Status = "invited"
	}

	// Try to save it in the database
	if err := env.Accounts.Insert(account); err != nil {
		utils.JSONResponse(w, 500, &AccountsCreateResponse{
			Success: false,
			Message: "Internal server error - AC/CR/02",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not insert an user to the database")
		return
	}

	// Send the email if classic and return a response
	if requestType == "classic" {
		// TODO: Send emails

		utils.JSONResponse(w, 201, &AccountsCreateResponse{
			Success: true,
			Message: "A new account was successfully created, you should receive a confirmation email soon™.",
			Account: account,
		})
		return
	}

	// Remove the token and return a response
	if requestType == "invited" {
		err := env.Tokens.DeleteID(input.Token)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err,
				"id":    input.Token,
			}).Error("Could not remove token from database")
		}

		utils.JSONResponse(w, 201, &AccountsCreateResponse{
			Success: true,
			Message: "A new account was successfully created",
			Account: account,
		})
		return
	}
}

// AccountsGetResponse contains the result of the AccountsGet request.
type AccountsGetResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	User    *models.Account `json:"user,omitempty"`
}

// AccountsGet returns the information about the specified account
func AccountsGet(c web.C, w http.ResponseWriter, r *http.Request) {
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
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err,
		}).Warn("Valid session referred to a removed account")

		// Try to remove the orphaned session
		if err := env.Tokens.DeleteID(session.ID); err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    session.ID,
				"error": err,
			}).Error("Unable to remove an orphaned session")
		} else if err := env.TokensCache.InvalidateToken(session.ID); err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    session.ID,
				"error": err,
			}).Error("Unable to remove an orphaned session from cache")
		} else {
			env.Log.WithFields(logrus.Fields{
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

// AccountsUpdateRequest contains the input for the AccountsUpdate endpoint.
type AccountsUpdateRequest struct {
	Type            string `json:"type" schema:"type"`
	AltEmail        string `json:"alt_email" schema:"alt_email"`
	CurrentPassword string `json:"current_password" schema:"current_password"`
	NewPassword     string `json:"new_password" schema:"new_password"`
}

// AccountsUpdateResponse contains the result of the AccountsUpdate request.
type AccountsUpdateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Account *models.Account `json:"account"`
}

// AccountsUpdate allows changing the account's information (password etc.)
func AccountsUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input AccountsUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 409, &AccountsUpdateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the account ID from the request
	id, ok := c.URLParams["id"]
	if !ok {
		utils.JSONResponse(w, 409, &AccountsUpdateResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsUpdateResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err,
		}).Warn("Valid session referred to a removed account")

		// Try to remove the orphaned session
		if err := env.Tokens.DeleteID(session.ID); err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    session.ID,
				"error": err,
			}).Error("Unable to remove an orphaned session")
		} else {
			env.Log.WithFields(logrus.Fields{
				"id": session.ID,
			}).Info("Removed an orphaned session")
		}

		utils.JSONResponse(w, 410, &AccountsUpdateResponse{
			Success: false,
			Message: "Account disabled",
		})
		return
	}

	if valid, _, err := user.VerifyPassword(input.CurrentPassword); err != nil || !valid {
		utils.JSONResponse(w, 409, &AccountsUpdateResponse{
			Success: false,
			Message: "Invalid current password",
		})
		return
	}

	err = user.SetPassword(input.NewPassword)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to hash a password")

		utils.JSONResponse(w, 500, &AccountsUpdateResponse{
			Success: false,
			Message: "Internal error (code AC/UP/01)",
		})
		return
	}

	if input.AltEmail != "" {
		user.AltEmail = input.AltEmail
	}

	err = env.Accounts.UpdateID(session.Owner, user)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to update an account")

		utils.JSONResponse(w, 500, &AccountsUpdateResponse{
			Success: false,
			Message: "Internal error (code AC/UP/02)",
		})
		return
	}

	utils.JSONResponse(w, 200, &AccountsUpdateResponse{
		Success: false,
		Message: "Your account has been successfully updated",
		Account: user,
	})
}

// AccountsDeleteResponse contains the result of the AccountsDelete request.
type AccountsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsDelete deletes an account and everything related to it.
func AccountsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the account ID from the request
	id, ok := c.URLParams["id"]
	if !ok {
		utils.JSONResponse(w, 409, &AccountsDeleteResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsDeleteResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err,
		}).Warn("Valid session referred to a removed account")

		// Try to remove the orphaned session
		if err := env.Tokens.DeleteID(session.ID); err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    session.ID,
				"error": err,
			}).Error("Unable to remove an orphaned session")
		} else {
			env.Log.WithFields(logrus.Fields{
				"id": session.ID,
			}).Info("Removed an orphaned session")
		}

		utils.JSONResponse(w, 410, &AccountsDeleteResponse{
			Success: false,
			Message: "Account disabled",
		})
		return
	}

	// TODO: Delete contacts

	// TODO: Delete emails

	// TODO: Delete labels

	// TODO: Delete threads

	// Delete tokens
	err = env.Tokens.DeleteByOwner(user.ID)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"id":    user.ID,
			"error": err,
		}).Error("Unable to remove account's tokens")

		utils.JSONResponse(w, 500, &AccountsDeleteResponse{
			Success: false,
			Message: "Internal error (code AC/DE/05)",
		})
		return
	}

	// Delete account
	err = env.Accounts.DeleteID(user.ID)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to delete an account")

		utils.JSONResponse(w, 500, &AccountsDeleteResponse{
			Success: false,
			Message: "Internal error (code AC/DE/06)",
		})
		return
	}

	utils.JSONResponse(w, 200, &AccountsDeleteResponse{
		Success: true,
		Message: "Your account has been successfully deleted",
	})
}

// AccountsWipeDataResponse contains the result of the AccountsWipeData request.
type AccountsWipeDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsWipeData wipes all data except the actual account and billing info.
func AccountsWipeData(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the account ID from the request
	id, ok := c.URLParams["id"]
	if !ok {
		utils.JSONResponse(w, 409, &AccountsWipeDataResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsWipeDataResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err,
		}).Warn("Valid session referred to a removed account")

		// Try to remove the orphaned session
		if err := env.Tokens.DeleteID(session.ID); err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    session.ID,
				"error": err,
			}).Error("Unable to remove an orphaned session")
		} else {
			env.Log.WithFields(logrus.Fields{
				"id": session.ID,
			}).Info("Removed an orphaned session")
		}

		utils.JSONResponse(w, 410, &AccountsWipeDataResponse{
			Success: false,
			Message: "Account disabled",
		})
		return
	}

	// TODO: Delete contacts

	// TODO: Delete emails

	// TODO: Delete labels

	// TODO: Delete threads

	// Delete tokens
	err = env.Tokens.DeleteByOwner(user.ID)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"id":    user.ID,
			"error": err,
		}).Error("Unable to remove account's tokens")

		utils.JSONResponse(w, 500, &AccountsWipeDataResponse{
			Success: false,
			Message: "Internal error (code AC/WD/05)",
		})
		return
	}

	utils.JSONResponse(w, 200, &AccountsWipeDataResponse{
		Success: true,
		Message: "Your account has been successfully wiped",
	})
}
