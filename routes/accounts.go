package routes

import (
	"encoding/json"
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
	Username   string `json:"username,omitempty" schema:"username"`
	Password   string `json:"password,omitempty" schema:"password"`
	AltEmail   string `json:"alt_email,omitempty" schema:"alt_email"`
	InviteCode string `json:"invite_code,omitempty" schema:"invite_code"`
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
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &AccountsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// TODO: Sanitize the username
	// TODO: Hash the password if it's not hashed already

	// Accounts flow:
	// 1) POST /accounts {username, alt_email}             => status = registered
	// 2) POST /accounts {username, invite_code}           => checks invite_code validity
	// 3) POST /accounts {username, invite_code, password} => status = setup
	requestType := "unknown"
	if input.Username != "" && input.Password == "" && input.AltEmail != "" && input.InviteCode == "" {
		requestType = "register"
	} else if input.Username != "" && input.Password == "" && input.AltEmail == "" && input.InviteCode != "" {
		requestType = "verify"
	} else if input.Username != "" && input.Password != "" && input.AltEmail == "" && input.InviteCode != "" {
		requestType = "setup"
	}

	// "unknown" requests are empty and invalid
	if requestType == "unknown" {
		utils.JSONResponse(w, 400, &AccountsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if requestType == "register" {
		// Normalize the username
		input.Username = utils.NormalizeUsername(input.Username)

		// Validate the username
		if len(input.Username) < 3 || len(input.Username) > 32 {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid username - it has to be at least 3 and at max 32 characters long",
			})
			return
		}

		// Ensure that the username is not used in address table
		if used, err := env.Addresses.GetAddress(utils.RemoveDots(input.Username)); err == nil || used != nil {
			utils.JSONResponse(w, 409, &AccountsCreateResponse{
				Success: false,
				Message: "Username already used",
			})
			return
		}

		// Then check it in the accounts table
		if ok, err := env.Accounts.IsUsernameUsed(utils.RemoveDots(input.Username)); ok || err != nil {
			utils.JSONResponse(w, 409, &AccountsCreateResponse{
				Success: false,
				Message: "Username already used",
			})
			return
		}

		// Also check that the email is unique
		if used, err := env.Accounts.IsEmailUsed(input.AltEmail); err != nil || used {
			if err != nil {
				env.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Unable to lookup registered accounts for emails")
			}

			utils.JSONResponse(w, 409, &AccountsCreateResponse{
				Success: false,
				Message: "Email already used",
			})
			return
		}

		// Both username and email are filled, so we can create a new account.
		account := &models.Account{
			Resource:   models.MakeResource("", utils.RemoveDots(input.Username)),
			StyledName: input.Username,
			Type:       "beta", // Is this the proper value?
			AltEmail:   input.AltEmail,
			Status:     "registered",
		}

		// Try to save it in the database
		if err := env.Accounts.Insert(account); err != nil {
			utils.JSONResponse(w, 500, &AccountsCreateResponse{
				Success: false,
				Message: "Internal server error - AC/CR/02",
			})

			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Could not insert an user into the database")
			return
		}

		// TODO: Send emails here. Depends on @andreis work.

		// Return information about the account
		utils.JSONResponse(w, 201, &AccountsCreateResponse{
			Success: true,
			Message: "Your account has been added to the beta queue",
			Account: account,
		})
		return
	} else if requestType == "verify" {
		// We're pretty much checking whether an invitation code can be used by the user
		input.Username = utils.RemoveDots(
			utils.NormalizeUsername(input.Username),
		)

		// Fetch the user from database
		account, err := env.Accounts.FindAccountByName(input.Username)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":    err.Error(),
				"username": input.Username,
			}).Warn("User not found in the database")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid username",
			})
			return
		}

		// Fetch the token from the database
		token, err := env.Tokens.GetToken(input.InviteCode)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Unable to fetch a registration token from the database")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation code",
			})
			return
		}

		// Ensure that the invite code was given to this particular user.
		if token.Owner != account.ID {
			env.Log.WithFields(logrus.Fields{
				"user_id": account.ID,
				"owner":   token.Owner,
			}).Warn("Not owned invitation code used by an user")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation code",
			})
			return
		}

		// Ensure that the token's type is valid
		if token.Type != "verify" {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation code",
			})
			return
		}

		// Check if it's expired
		if token.Expired() {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Expired invitation code",
			})
			return
		}

		// Ensure that the account is "registered"
		if account.Status != "registered" {
			utils.JSONResponse(w, 403, &AccountsCreateResponse{
				Success: true,
				Message: "This account was already configured",
			})
			return
		}

		// Everything is fine, return it.
		utils.JSONResponse(w, 200, &AccountsCreateResponse{
			Success: true,
			Message: "Valid token was provided",
		})
		return
	} else if requestType == "setup" {
		// User is setting the password in the setup wizard. This should be one of the first steps,
		// as it's required for him to acquire an authentication token to configure their account.
		input.Username = utils.RemoveDots(
			utils.NormalizeUsername(input.Username),
		)

		// Fetch the user from database
		account, err := env.Accounts.FindAccountByName(input.Username)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":    err.Error(),
				"username": input.Username,
			}).Warn("User not found in the database")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid username",
			})
			return
		}

		// Fetch the token from the database
		token, err := env.Tokens.GetToken(input.InviteCode)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Unable to fetch a registration token from the database")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation code",
			})
			return
		}

		// Ensure that the invite code was given to this particular user.
		if token.Owner != account.ID {
			env.Log.WithFields(logrus.Fields{
				"user_id": account.ID,
				"owner":   token.Owner,
			}).Warn("Not owned invitation code used by an user")

			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation code",
			})
			return
		}

		// Ensure that the token's type is valid
		if token.Type != "verify" {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Invalid invitation code",
			})
			return
		}

		// Check if it's expired
		if token.Expired() {
			utils.JSONResponse(w, 400, &AccountsCreateResponse{
				Success: false,
				Message: "Expired invitation code",
			})
			return
		}

		// Ensure that the account is "registered"
		if account.Status != "registered" {
			utils.JSONResponse(w, 403, &AccountsCreateResponse{
				Success: true,
				Message: "This account was already configured",
			})
			return
		}

		// Our token is fine, next part: password.

		// Ensure that user has chosen a secure password (check against 10k most used)
		if !utils.IsPasswordSecure(input.Password) {
			utils.JSONResponse(w, 403, &AccountsCreateResponse{
				Success: false,
				Message: "Weak password",
			})
			return
		}

		// We can't really make more checks on the password, user could as well send us a hash
		// of a simple password, but we assume that no developer is that stupid (actually,
		// considering how many people upload their private keys and AWS credentials, I'm starting
		// to doubt the competence of some so-called "web deyvelopayrs")

		// Set the password
		err = account.SetPassword(input.Password)
		if err != nil {
			utils.JSONResponse(w, 500, &AccountsCreateResponse{
				Success: false,
				Message: "Internal server error - AC/CR/01",
			})

			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to hash the password")
			return
		}

		account.Status = "setup"

		// Create labels
		err = env.Labels.Insert([]*models.Label{
			&models.Label{
				Resource: models.MakeResource(account.ID, "Inbox"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Sent"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Drafts"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Trash"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Spam"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Starred"),
				Builtin:  true,
			},
		})
		if err != nil {
			utils.JSONResponse(w, 500, &AccountsCreateResponse{
				Success: false,
				Message: "Internal server error - AC/CR/03",
			})

			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Could not insert labels into the database")
			return
		}

		// Add a new mapping
		err = env.Addresses.Insert(&models.Address{
			Resource: models.Resource{
				ID:           account.Name,
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				Owner:        account.ID,
			},
		})
		if err != nil {
			utils.JSONResponse(w, 500, &AccountsCreateResponse{
				Success: false,
				Message: "Unable to create a new address mapping",
			})

			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Could not insert an address mapping into db")
			return
		}

		// Update the account
		err = env.Accounts.UpdateID(account.ID, account)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    account.ID,
			}).Error("Unable to update an account")

			utils.JSONResponse(w, 500, &AccountsCreateResponse{
				Success: false,
				Message: "Unable to update the account",
			})
			return
		}

		// Remove the token and return a response
		err = env.Tokens.DeleteID(input.InviteCode)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    input.InviteCode,
			}).Error("Could not remove the token from database")
		}

		utils.JSONResponse(w, 200, &AccountsCreateResponse{
			Success: true,
			Message: "Your account has been initialized successfully",
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
	AltEmail        string      `json:"alt_email" schema:"alt_email"`
	CurrentPassword string      `json:"current_password" schema:"current_password"`
	NewPassword     string      `json:"new_password" schema:"new_password"`
	FactorType      string      `json:"factor_type" schema:"factor_type"`
	FactorValue     []string    `json:"factor_value" schema:"factor_value"`
	Token           string      `json:"token" schema:"token"`
	Settings        interface{} `json:"settings" schema:"settings"`
	PublicKey       string      `json:"public_key" schema:"public_key"`
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
			"error": err.Error(),
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

	if input.NewPassword != "" {
		if valid, _, err := user.VerifyPassword(input.CurrentPassword); err != nil || !valid {
			utils.JSONResponse(w, 403, &AccountsUpdateResponse{
				Success: false,
				Message: "Invalid current password",
			})
			return
		}
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
				"error": err.Error(),
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

	if input.Settings != nil {
		user.Settings = input.Settings
	}

	if input.PublicKey != "" {
		key, err := env.Keys.FindByFingerprint(input.PublicKey)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": input.PublicKey,
			}).Error("Unable to find a key")

			utils.JSONResponse(w, 400, &AccountsUpdateResponse{
				Success: false,
				Message: "Invalid public key",
			})
			return
		}

		if key.Owner != user.ID {
			env.Log.WithFields(logrus.Fields{
				"user_id:":    user.ID,
				"owner":       key.Owner,
				"fingerprint": input.PublicKey,
			}).Error("Unable to find a key")

			utils.JSONResponse(w, 400, &AccountsUpdateResponse{
				Success: false,
				Message: "Invalid public key",
			})
			return
		}

		user.PublicKey = input.PublicKey
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
			"error": err.Error(),
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
			"error": err.Error(),
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
			"error": err.Error(),
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
			"error": err.Error(),
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
			"error": err.Error(),
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

type AccountsStartOnboardingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func AccountsStartOnboarding(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the account ID from the request
	id := c.URLParams["id"]

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, &AccountsStartOnboardingResponse{
			Success: false,
			Message: `Only the "me" user is implemented`,
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	account, err := env.Accounts.GetTokenOwner(session)
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err.Error(),
		}).Warn("Valid session referred to a removed account")

		utils.JSONResponse(w, 410, &AccountsStartOnboardingResponse{
			Success: false,
			Message: "Account disabled",
		})
		return
	}

	x1, ok := account.Settings.(map[string]interface{})
	if !ok {
		utils.JSONResponse(w, 403, &AccountsStartOnboardingResponse{
			Success: false,
			Message: "Account misconfigured #1",
		})
		return
	}

	x2, ok := x1["firstName"]
	if !ok {
		utils.JSONResponse(w, 403, &AccountsStartOnboardingResponse{
			Success: false,
			Message: "Account misconfigured #2",
		})
		return
	}

	x3, ok := x2.(string)
	if !ok {
		utils.JSONResponse(w, 403, &AccountsStartOnboardingResponse{
			Success: false,
			Message: "Account misconfigured #3",
		})
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"type":  "onboarding",
		"email": account.Name + "@lavaboom.com",
		// polish roulette
		"first_name": x3,
	})
	if !ok {
		utils.JSONResponse(w, 500, &AccountsStartOnboardingResponse{
			Success: false,
			Message: "Unable to encode a message",
		})
		return
	}

	if err := env.Producer.Publish("hub", data); err != nil {
		utils.JSONResponse(w, 500, &AccountsCreateResponse{
			Success: false,
			Message: "Unable to initialize onboarding emails",
		})
		return
	}

	utils.JSONResponse(w, 200, &AccountsStartOnboardingResponse{
		Success: true,
		Message: "Onboarding emails for your account have been initialized",
	})
}
