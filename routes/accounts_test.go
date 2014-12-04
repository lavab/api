package routes_test

import (
	"testing"
	"time"

	"github.com/franela/goreq"
	"github.com/stretchr/testify/require"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
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
	env.Log.Print(response)
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

func TestAccountsCreateInvited(t *testing.T) {
	const (
		username = "jeremy"
		password = "potato"
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
	result1, err := goreq.Request{
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
	var response1 routes.AccountsCreateResponse
	err = result1.Body.FromJsonTo(&response1)
	require.Nil(t, err)

	// Check the result's contents
	require.True(t, response1.Success)
	require.Equal(t, "A new account was successfully created", response1.Message)
	require.NotEmpty(t, response1.Account.ID)

	accountID = response1.Account.ID

	// POST /accounts - invited with wrong token
	result2, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username + "2",
			Password: password,
			Token:    "asdasdasd",
		},
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var response2 routes.AccountsCreateResponse
	err = result2.Body.FromJsonTo(&response2)
	require.Nil(t, err)

	// Check the result's contents
	require.False(t, response2.Success)
	require.Equal(t, "Invalid invitation token", response2.Message)
}

func TestAccountsCreateInvitedExisting(t *testing.T) {
	const (
		username = "jeremy"
		password = "potato"
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
	require.Equal(t, "Username already exists", response.Message)
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
		username = "jeremy_was_invited"
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

func TestAccountsPrepareToken(t *testing.T) {
	// log in as mr jeremy potato
	const (
		username = "jeremy"
		password = "potato"
	)
	// POST /accounts - classic
	request, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/tokens",
		ContentType: "application/json",
		Body: routes.TokensCreateRequest{
			Username: username,
			Password: password,
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
