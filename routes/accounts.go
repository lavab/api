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

// AccountsListResponse contains the result of the AccountsList request.
type AccountsListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountsList returns a list of accounts visible to an user
func AccountsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &AccountsListResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// AccountsCreateRequest contains the input for the AccountsCreate endpoint.
type AccountsCreateRequest struct {
	Token    string `json:"token,omitempty" schema:"token"`
	Username string `json:"username,omitempty" schema:"username"`
	Password string `json:"password,omitempty" schema:"password"`
	AltEmail string `json:"alt_email,omitempty" schema:"alt_email"`
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
	// 4) alt_email + username            - register for beta with username reservation
	requestType := "unknown"
	if input.AltEmail == "" && input.Username != "" && input.Password != "" && input.Token != "" {
		requestType = "invited"
	} else if input.AltEmail != "" && input.Username != "" && input.Password != "" && input.Token == "" {
		requestType = "classic"
	} else if input.AltEmail != "" && input.Username == "" && input.Password == "" && input.Token == "" {
		requestType = "queue/classic"
	} else if input.AltEmail != "" && input.Username != "" && input.Password == "" && input.Token == "" {
		requestType = "queue/reserve"
	}

	// "unknown" requests are empty and invalid
	if requestType == "unknown" {
		utils.JSONResponse(w, 400, &AccountsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if input.Username != "" {
		if used, err := env.Reservations.IsUsernameUsed(input.Username); err != nil || used {
			if err != nil {
				env.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Unable to lookup reservations for usernames")
			}

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Username already reserved",
			})
			return
		}

		if used, err := env.Accounts.IsUsernameUsed(input.Username); err != nil || used {
			if err != nil {
				env.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Unable to lookup registered accounts for usernames")
			}

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Username already used",
			})
			return
		}
	}

	// Adding to [beta] queue
	if requestType[:5] == "queue" {
		if requestType[6:] == "reserve" {
			// Is username reservation enabled?
			if !env.Config.UsernameReservation {
				utils.JSONResponse(w, 403, &AccountsCreateResponse{
					Success: false,
					Message: "Username reservation is disabled",
				})
				return
			}
		}

		// Ensure that the email is not already used to reserve/register
		if used, err := env.Reservations.IsEmailUsed(input.AltEmail); err != nil || used {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Email already used for a reservation",
			})
			return
		}

		if used, err := env.Accounts.IsEmailUsed(input.AltEmail); err != nil || used {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Email already used for a reservation",
			})
			return
		}

		// Prepare data to insert
		reservation := &models.Reservation{
			Email:    input.AltEmail,
			Resource: models.MakeResource("", input.Username),
		}

		err := env.Reservations.Insert(reservation)
		if err != nil {
			utils.JSONResponse(w, 500, &AccountsCreateResponse{
				Success: false,
				Message: "Internal error while reserving the account",
			})
			return
		}

		utils.JSONResponse(w, 201, &AccountsCreateResponse{
			Success: true,
			Message: "Reserved an account",
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

	// Check for generic passwords
	if input.Password != "" && !utils.IsPasswordSecure(input.Password) {
		utils.JSONResponse(w, 403, &AccountsCreateResponse{
			Success: false,
			Message: "Weak password",
		})
		return
	}

	// Both username and password are filled, so we can create a new account.
	account := &models.Account{
		Resource: models.MakeResource("", input.Username),
		Type:     "beta",
		AltEmail: input.AltEmail,
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

		account.AltEmail = token.Email
	}

	// TODO: sanitize user name (i.e. remove caps, periods)

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
		}).Error("Could not insert an user into the database")
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
	Account *models.Account `json:"user,omitempty"`
}

// AccountsGet returns the information about the specified account
func AccountsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the account ID from the request
	id := c.URLParams["id"]

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsGetResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, &AccountsDeleteResponse{
			Success: false,
			Message: "Unable to resolve the account",
		})
		return
	}

	// Return the user struct
	utils.JSONResponse(w, 200, &AccountsGetResponse{
		Success: true,
		Account: user,
	})
}

// AccountsUpdateRequest contains the input for the AccountsUpdate endpoint.
type AccountsUpdateRequest struct {
	AltEmail        string   `json:"alt_email" schema:"alt_email"`
	CurrentPassword string   `json:"current_password" schema:"current_password"`
	NewPassword     string   `json:"new_password" schema:"new_password"`
	FactorType      string   `json:"factor_type" schema:"factor_type"`
	FactorValue     []string `json:"factor_value" schema:"factor_value"`
	Token           string   `json:"token" schema:"token"`
}

// AccountsUpdateResponse contains the result of the AccountsUpdate request.
type AccountsUpdateResponse struct {
	Success         bool            `json:"success"`
	Message         string          `json:"message,omitempty"`
	Account         *models.Account `json:"account,omitempty"`
	FactorType      string          `json:"factor_type,omitempty"`
	FactorChallenge string          `json:"factor_challenge,omitempty"`
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

		utils.JSONResponse(w, 400, &AccountsUpdateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the account ID from the request
	id := c.URLParams["id"]

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsUpdateResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, &AccountsDeleteResponse{
			Success: false,
			Message: "Unable to resolve the account",
		})
		return
	}

	if valid, _, err := user.VerifyPassword(input.CurrentPassword); err != nil || !valid {
		utils.JSONResponse(w, 403, &AccountsUpdateResponse{
			Success: false,
			Message: "Invalid current password",
		})
		return
	}

	// Check for 2nd factor
	if user.FactorType != "" {
		factor, ok := env.Factors[user.FactorType]
		if ok {
			// Verify the 2FA
			verified, challenge, err := user.Verify2FA(factor, input.Token)
			if err != nil {
				utils.JSONResponse(w, 500, &AccountsUpdateResponse{
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
				utils.JSONResponse(w, 403, &AccountsUpdateResponse{
					Success:         false,
					Message:         "2FA token was not passed",
					FactorType:      user.FactorType,
					FactorChallenge: challenge,
				})
				return
			}

			// Token was incorrect
			if !verified {
				utils.JSONResponse(w, 403, &AccountsUpdateResponse{
					Success:    false,
					Message:    "Invalid token passed",
					FactorType: user.FactorType,
				})
				return
			}
		}
	}

	if input.NewPassword != "" && !utils.IsPasswordSecure(input.NewPassword) {
		utils.JSONResponse(w, 400, &AccountsUpdateResponse{
			Success: false,
			Message: "Weak new password",
		})
		return
	}

	if input.NewPassword != "" {
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
	}

	if input.AltEmail != "" {
		user.AltEmail = input.AltEmail
	}

	if input.FactorType != "" {
		// Check if such factor exists
		if _, exists := env.Factors[input.FactorType]; !exists {
			utils.JSONResponse(w, 400, &AccountsUpdateResponse{
				Success: false,
				Message: "Invalid new 2FA type",
			})
			return
		}

		user.FactorType = input.FactorType
	}

	if input.FactorValue != nil && len(input.FactorValue) > 0 {
		user.FactorValue = input.FactorValue
	}

	user.DateModified = time.Now()

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
		Success: true,
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
	id := c.URLParams["id"]

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsDeleteResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, &AccountsDeleteResponse{
			Success: false,
			Message: "Unable to resolve the account",
		})
		return
	}

	// TODO: Delete contacts

	// TODO: Delete emails

	// TODO: Delete labels

	// TODO: Delete threads

	// Delete tokens
	err = env.Tokens.DeleteOwnedBy(user.ID)
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
	id := c.URLParams["id"]

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsWipeDataResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetTokenOwner(session)
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err,
		}).Warn("Valid session referred to a removed account")

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
	err = env.Tokens.DeleteOwnedBy(user.ID)
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
