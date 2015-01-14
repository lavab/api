package routes_test

import (
	"testing"
	//"time"

	"github.com/franela/goreq"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	//"github.com/lavab/api/env"
	//"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
)

var (
	accountUsername string
	accountPassword string
	accountID       string
)

func TestAccountsCreateInvalid(t *testing.T) {
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body:        "!@#!@#",
	}.Do()
	require.Nil(t, err)

	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)
	require.False(t, response.Success)
	require.Equal(t, "Invalid input format", response.Message)
}

func TestAccountsCreateUnknown(t *testing.T) {
	// POST /accounts - unknown
	result, err := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts",
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check values
	require.False(t, response.Success)
	require.Equal(t, "Invalid request", response.Message)
}

func TestAccountsCreateRegister(t *testing.T) {
	const (
		username = "jeremy"
		password = "potato"
		email    = "jeremy@potato.org"
	)

	// Prepare account information
	accountUsername = username
	passwordHash := sha3.Sum256([]byte(password))
	accountPassword = string(passwordHash[:])

	// Prepare a token
	/*inviteToken := models.Token{
		Resource: models.MakeResource("", "test invite token"),
		Type:     "invite",
	}
	inviteToken.ExpireSoon()

	err := env.Tokens.Insert(inviteToken)
	require.Nil(t, err)*/

	// POST /accounts - invited
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username,
			AltEmail: email,
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Your account has been added to the beta queue", response.Message)
	require.True(t, response.Success)
	require.NotEmpty(t, response.Account.ID)

	accountID = response.Account.ID
}

func TestAccountsCreateInvitedExistingUsername(t *testing.T) {
	const (
		username = "jeremy"
		email    = "jeremy@potato.org"
	)

	// POST /accounts - invited
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username,
			AltEmail: email,
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Username already used", response.Message)
}

func TestAccountsCreateInvitedExistingEmail(t *testing.T) {
	const (
		username = "jeremy2"
		email    = "jeremy@potato.org"
	)

	// POST /accounts - invited
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username,
			AltEmail: email,
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Email already used", response.Message)
}

/*
func TestAccountsCreateInvitedWeakPassword(t *testing.T) {
	const (
		username = "jeremylicious"
		password = "c0067d4af4e87f00dbac63b6156828237059172d1bbeac67427345d6a9fda484"
	)

	// Prepare a token
	inviteToken := models.Token{
		Resource: models.MakeResource("", "test invite token"),
		Type:     "invite",
	}
	inviteToken.ExpireSoon()

	err := env.Tokens.Insert(inviteToken)
	require.Nil(t, err)

	// POST /accounts - invited
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username,
			Password: password,
			Token:    inviteToken.ID,
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Weak password", response.Message)
}

func TestAccountsCreateInvitedExpired(t *testing.T) {
	const (
		username = "jeremy2"
		password = "potato2"
	)

	// Prepare a token
	inviteToken := models.Token{
		Resource: models.MakeResource("", "test invite token"),
		Type:     "invite",
	}
	inviteToken.ExpiryDate = time.Now().Truncate(time.Hour)

	err := env.Tokens.Insert(inviteToken)
	require.Nil(t, err)

	// POST /accounts - invited
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username,
			Password: password,
			Token:    inviteToken.ID,
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Expired invitation token", response.Message)
}

func TestAccountsCreateInvitedWrongType(t *testing.T) {
	const (
		username = "jeremy2"
		password = "potato2"
	)

	// Prepare a token
	inviteToken := models.Token{
		Resource: models.MakeResource("", "test not invite token"),
		Type:     "not invite",
	}
	inviteToken.ExpiryDate = time.Now().Truncate(time.Hour)

	err := env.Tokens.Insert(inviteToken)
	require.Nil(t, err)

	// POST /accounts - invited
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username,
			Password: password,
			Token:    inviteToken.ID,
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Invalid invitation token", response.Message)
}

func TestAccountsCreateClassic(t *testing.T) {
	const (
		username = "jeremy"
		password = "potato"
	)

	// POST /accounts - classic
	createClassicResult, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username + "classic",
			Password: password,
			AltEmail: "something@example.com",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var createClassicResponse routes.AccountsCreateResponse
	err = createClassicResult.Body.FromJsonTo(&createClassicResponse)
	require.Nil(t, err)

	// Check the result's contents
	require.True(t, createClassicResponse.Success)
	require.Equal(t, "A new account was successfully created, you should receive a confirmation email soon™.", createClassicResponse.Message)
	require.NotEmpty(t, createClassicResponse.Account.ID)
}

func TestAccountsCreateClassicDisabled(t *testing.T) {
	const (
		username = "jeremy_was_invited"
		password = "potato"
	)

	env.Config.ClassicRegistration = false

	// POST /accounts - classic
	createClassicResult, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username + "classic",
			Password: password,
			AltEmail: "something@example.com",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var createClassicResponse routes.AccountsCreateResponse
	err = createClassicResult.Body.FromJsonTo(&createClassicResponse)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, createClassicResponse.Success)
	require.Equal(t, "Classic registration is disabled", createClassicResponse.Message)

	env.Config.ClassicRegistration = true
}

func TestAccountsCreateQueueReservation(t *testing.T) {
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "reserved@example.com",
			Username: "reserved",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Reserved an account", response.Message)
	require.True(t, response.Success)
}

func TestAccountsCreateQueueReservationUsernameReserved(t *testing.T) {
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "not-reserved@example.com",
			Username: "reserved",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Username already reserved", response.Message)
	require.False(t, response.Success)
}

func TestAccountsCreateQueueReservationUsernameTaken(t *testing.T) {
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "not-reserved@example.com",
			Username: "jeremy",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Username already used", response.Message)
}

func TestAccountsCreateQueueReservationDisabled(t *testing.T) {
	env.Config.UsernameReservation = false
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "something@example.com",
			Username: "something",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Username reservation is disabled", response.Message)
	env.Config.UsernameReservation = true
}

func TestAccountsCreateQueueClassicUsedEmail(t *testing.T) {
	// POST /accounts - queue
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "something@example.com",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Email already used for a reservation", response.Message)
	require.False(t, response.Success)
}

func TestAccountsCreateQueueClassicReservedEmail(t *testing.T) {
	// POST /accounts - queue
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "reserved@example.com",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Email already used for a reservation", response.Message)
	require.False(t, response.Success)
}


func TestAccountsPrepareToken(t *testing.T) {
	// POST /accounts - classic
	request, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/tokens",
		ContentType: "application/json",
		Body: routes.TokensCreateRequest{
			Username: accountUsername,
			Password: accountPassword,
			Type:     "auth",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.TokensCreateResponse
	err = request.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.True(t, response.Success)
	require.Equal(t, "Authentication successful", response.Message)
	require.NotEmpty(t, response.Token.ID)

	// Populate the global token variable
	authToken = response.Token.ID
}

func TestAccountsList(t *testing.T) {
	// GET /accounts
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsListResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, "Sorry, not implemented yet", response.Message)
}

func TestAccountsGetMe(t *testing.T) {
	// GET /accounts/me
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.True(t, response.Success)
	require.Equal(t, "jeremy", response.Account.Name)
}

func TestAccountsGetNotMe(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/not-me",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response.Success)
	require.Equal(t, `Only the "me" user is implemented`, response.Message)
}

func TestAccountUpdateMe(t *testing.T) {
	// PUT /accounts/me
	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/accounts/me",
		ContentType: "application/json",
		Body: &routes.AccountsUpdateRequest{
			CurrentPassword: "potato",
			NewPassword:     "cabbage",
			AltEmail:        "john.cabbage@example.com",
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Your account has been successfully updated", response.Message)
	require.True(t, response.Success)
	require.Equal(t, "jeremy", response.Account.Name)
	require.Equal(t, "john.cabbage@example.com", response.Account.AltEmail)
}

func TestAccountUpdateInvalid(t *testing.T) {
	// PUT /accounts/me
	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/accounts/me",
		ContentType: "application/json",
		Body:        "123123123!@#!@#!@#",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Invalid input format", response.Message)
	require.False(t, response.Success)
}

func TestAccountUpdateNotMe(t *testing.T) {
	// PUT /accounts/me
	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/accounts/not-me",
		ContentType: "application/json",
		Body: &routes.AccountsUpdateRequest{
			CurrentPassword: "potato",
			NewPassword:     "cabbage",
			AltEmail:        "john.cabbage@example.com",
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, `Only the "me" user is implemented`, response.Message)
	require.False(t, response.Success)
}

func TestAccountUpdateMeInvalidPassword(t *testing.T) {
	// PUT /accounts/me
	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/accounts/me",
		ContentType: "application/json",
		Body: &routes.AccountsUpdateRequest{
			CurrentPassword: "potato2",
			NewPassword:     "cabbage",
			AltEmail:        "john.cabbage@example.com",
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Invalid current password", response.Message)
	require.False(t, response.Success)
}

func TestAccountsWipeDataNotMe(t *testing.T) {
	// POST /accounts/me/wipe-data
	request := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts/not-me/wipe-data",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, `Only the "me" user is implemented`, response.Message)
	require.False(t, response.Success)
}

func TestAccountsWipeData(t *testing.T) {
	// POST /accounts/me/wipe-data
	request := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts/me/wipe-data",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Your account has been successfully wiped", response.Message)
	require.True(t, response.Success)
}

func TestAccountsDeleteNotMe(t *testing.T) {
	// Prepare a token
	token := models.Token{
		Resource: models.MakeResource(accountID, "test invite token"),
		Type:     "auth",
	}
	token.ExpireSoon()

	err := env.Tokens.Insert(token)
	require.Nil(t, err)

	// DELETE /accounts/me
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/accounts/not-me",
	}
	request.AddHeader("Authorization", "Bearer "+token.ID)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, `Only the "me" user is implemented`, response.Message)
	require.False(t, response.Success)
}

func TestAccountsDelete(t *testing.T) {
	// Prepare a token
	token := models.Token{
		Resource: models.MakeResource(accountID, "test invite token"),
		Type:     "auth",
	}
	token.ExpireSoon()

	err := env.Tokens.Insert(token)
	require.Nil(t, err)

	// DELETE /accounts/me
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "Bearer "+token.ID)
	result, err := request.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	// Check the result's contents
	require.Equal(t, "Your account has been successfully deleted", response.Message)
	require.True(t, response.Success)
}
*/
