package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// AccountsList returns a list of accounts visible to an user
func AccountsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, utils.NewError(
		utils.AccountsListUnknown, "Account not implemented yet", false,
	))
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.AccountsCreateInvalidInput, err, false,
		))
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.AccountsCreateUnknownStep, "Unable to recognize the step", false,
		))
		return
	}

	if requestType == "register" {
		// Normalize the username
		input.Username = utils.NormalizeUsername(input.Username)

		// Validate the username
		if len(input.Username) < 3 || len(utils.RemoveDots(input.Username)) < 3 || len(input.Username) > 32 {
			utils.JSONResponse(w, 400, utils.NewError(
				utils.AccountsCreateInvalidLength, "Invalid username - it has to be at least 3 and at max 32 characters long", false,
			))
			return
		}

		// Ensure that the username is not used in address table
		if used, err := env.Addresses.GetAddress(utils.RemoveDots(input.Username)); err == nil || used != nil {
			utils.JSONResponse(w, 409, utils.NewError(
				utils.AccountsCreateUsernameTaken, "Username already taken", false,
			))
			return
		}

		// Then check it in the accounts table
		if ok, err := env.Accounts.IsUsernameUsed(utils.RemoveDots(input.Username)); ok || err != nil {
			utils.JSONResponse(w, 409, utils.NewError(
				utils.AccountsCreateUsernameTaken, "Username already taken", false,
			))
			return
		}

		// Also check that the email is unique
		if used, err := env.Accounts.IsEmailUsed(input.AltEmail); err != nil || used {
			utils.JSONResponse(w, 409, utils.NewError(
				utils.AccountsCreateEmailUsed, "Email already used", false,
			))
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
			utils.JSONResponse(w, 500, utils.NewError(
				utils.AccountsCreateUnableToInsertAccount, err, true,
			))
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
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateUserNotFound, err, false,
			))
			return
		}

		// Fetch the token from the database
		token, err := env.Tokens.GetToken(input.InviteCode)
		if err != nil {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateInvalidToken, err, false,
			))
			return
		}

		// Ensure that the invite code was given to this particular user.
		if token.Owner != account.ID {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateInvalidTokenOwner, "You don't own this invitation code", false,
			))
			return
		}

		// Ensure that the token's type is valid
		if token.Type != "verify" {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateInvalidTokenType, "Invalid token type - "+token.Type, false,
			))
			return
		}

		// Check if it's expired
		if token.Expired() {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateExpiredToken, "This token has expired", false,
			))
			return
		}

		// Ensure that the account is "registered"
		if account.Status != "registered" {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateAlreadyConfigured, "This account is already configured", false,
			))
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
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateUserNotFound, err, false,
			))
			return
		}

		// Fetch the token from the database
		token, err := env.Tokens.GetToken(input.InviteCode)
		if err != nil {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateInvalidToken, err, false,
			))
			return
		}

		// Ensure that the invite code was given to this particular user.
		if token.Owner != account.ID {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateInvalidTokenOwner, "You don't own this invitation code", false,
			))
			return
		}

		// Ensure that the token's type is valid
		if token.Type != "verify" {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateInvalidTokenType, "Invalid token type - "+token.Type, false,
			))
			return
		}

		// Check if it's expired
		if token.Expired() {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateExpiredToken, "This token has expired", false,
			))
			return
		}

		// Ensure that the account is "registered"
		if account.Status != "registered" {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateAlreadyConfigured, "This account is already configured", false,
			))
			return
		}

		// Our token is fine, next part: password.

		// Ensure that user has chosen a secure password (check against 10k most used)
		if env.PasswordBF.TestString(input.Password) {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsCreateWeakPassword, "Weak password", false,
			))
			return
		}

		// We can't really make more checks on the password, user could as well send us a hash
		// of a simple password, but we assume that no developer is that stupid (actually,
		// considering how many people upload their private keys and AWS credentials to GitHub,
		// I'm starting to doubt the competence of some so-called "web deyvelopayrs")

		// Set the password
		err = account.SetPassword(input.Password)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.AccountsCreateUnableToHash, err, true,
			))
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
			utils.JSONResponse(w, 500, utils.NewError(
				utils.AccountsCreateUnableToPrepareLabels, err, true,
			))
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
			utils.JSONResponse(w, 500, utils.NewError(
				utils.AccountsCreateUnableToCreateAddress, err, true,
			))
			return
		}

		// Update the account
		err = env.Accounts.UpdateID(account.ID, account)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.AccountsCreateUnableToUpdateAccount, err, true,
			))
			return
		}

		// Remove the token and return a response
		err = env.Tokens.DeleteID(input.InviteCode)
		if err != nil {
			env.Raven.CaptureError(err, nil)
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
		utils.JSONResponse(w, 501, utils.NewError(
			utils.AccountsGetOnlyMe, "You can only get your own account's details", false,
		))
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsGetUnableToGet, err, true,
		))
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
	Token           string      `json:"token" schema:"token"`
	Settings        interface{} `json:"settings" schema:"settings"`
	PublicKey       string      `json:"public_key" schema:"public_key"`
}

// AccountsUpdateResponse contains the result of the AccountsUpdate request.
type AccountsUpdateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Account *models.Account `json:"account,omitempty"`
}

// AccountsUpdate allows changing the account's information (password etc.)
func AccountsUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input AccountsUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.AccountsUpdateOnlyMe, err, false,
		))
		return
	}

	// Get the account ID from the request
	id := c.URLParams["id"]

	// Right now we only support "me" as the ID
	if id != "me" {
		utils.JSONResponse(w, 501, utils.NewError(
			utils.AccountsUpdateOnlyMe, "You can only update your own account", false,
		))
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsUpdateUnableToGet, err, true,
		))
		return
	}

	if input.NewPassword != "" {
		if valid, _, err := user.VerifyPassword(input.CurrentPassword); err != nil || !valid {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsUpdateInvalidCurrentPassword, err, false,
			))
			return
		}
	}

	if input.NewPassword != "" && env.PasswordBF.TestString(input.NewPassword) {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.AccountsUpdateWeakPassword, "Weak new password", false,
		))
		return
	}

	if input.NewPassword != "" {
		err = user.SetPassword(input.NewPassword)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.AccountsUpdateUnableToHash, err, true,
			))
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
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsUpdateInvalidPublicKey, err, false,
			))
			return
		}

		if key.Owner != user.ID {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.AccountsUpdateInvalidPublicKeyOwner, "You're not the owner of that public key", false,
			))
			return
		}

		user.PublicKey = input.PublicKey
	}

	user.DateModified = time.Now()

	err = env.Accounts.UpdateID(session.Owner, user)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsUpdateUnableToUpdate, err, true,
		))
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
		utils.JSONResponse(w, 501, utils.NewError(
			utils.AccountsDeleteOnlyMe, "You can only delete your own account", false,
		))
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsDeleteUnableToGet, err, true,
		))
		return
	}

	// TODO: Delete contacts

	// TODO: Delete emails

	// TODO: Delete labels

	// TODO: Delete threads

	// Delete tokens
	err = env.Tokens.DeleteOwnedBy(user.ID)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsDeleteUnableToDelete, err, true,
		))
		return
	}

	// Delete account
	err = env.Accounts.DeleteID(user.ID)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsDeleteUnableToDelete, err, true,
		))
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
		utils.JSONResponse(w, 501, utils.NewError(
			utils.AccountsWipeDataOnlyMe, "You can only delete your own account", false,
		))
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	user, err := env.Accounts.GetTokenOwner(session)
	if err != nil {
		// The session refers to a non-existing user
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsWipeDataUnableToGet, err, true,
		))
		return
	}

	// TODO: Delete contacts

	// TODO: Delete emails

	// TODO: Delete labels

	// TODO: Delete threads

	// Delete tokens
	err = env.Tokens.DeleteOwnedBy(user.ID)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsWipeDataUnableToDelete, err, true,
		))
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
		utils.JSONResponse(w, 501, utils.NewError(
			utils.AccountsStartOnboardingOnlyMe, "You can only start onboarding for your own account", false,
		))
		return
	}

	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Fetch the user object from the database
	account, err := env.Accounts.GetTokenOwner(session)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsStartOnboardingUnableToGet, err, true,
		))
		return
	}

	x1, ok := account.Settings.(map[string]interface{})
	if !ok {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.AccountsStartOnboardingMisconfigured, "Account settings are not an array", true,
		))
		return
	}

	x2, ok := x1["firstName"]
	if !ok {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.AccountsStartOnboardingMisconfigured, "Account settings do not have a first name property", true,
		))
		return
	}

	x3, ok := x2.(string)
	if !ok {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.AccountsStartOnboardingMisconfigured, "First name in account settings is not a string", true,
		))
		return
	}

	data, _ := json.Marshal(map[string]interface{}{
		"type":       "onboarding",
		"email":      account.Name + "@lavaboom.com",
		"first_name": x3,
	})

	if err := env.Producer.Publish("hub", data); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.AccountsStartOnboardingUnableToInit, err, true,
		))
		return
	}

	utils.JSONResponse(w, 200, &AccountsStartOnboardingResponse{
		Success: true,
		Message: "Onboarding emails for your account have been initialized",
	})
}
