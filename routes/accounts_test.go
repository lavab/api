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
	// POST /accounts - invalid
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body:        "!@#!@#",
	}.Do()
	require.Nil(t, err, "querying invalid /accounts should not return an error")

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	env.Log.Print(response)
	require.Nil(t, err, "unmarshaling invalid account creation should not return an error")
	require.False(t, response.Success, "request should fail")
	require.Equal(t, "Invalid input format", response.Message, "proper message should be returned")
}

func TestAccountsCreateUnknown(t *testing.T) {
	// POST /accounts - unknown
	result, err := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts",
	}.Do()
	require.Nil(t, err, "querying unknown /accounts should not return an error")

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check values
	require.False(t, response.Success, "unknown request should return success false")
	require.Equal(t, "Invalid request", response.Message, "unknown request should return proper error msg")
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
	require.Nil(t, err, "inserting a new invitation token should not return an error")

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
	require.Nil(t, err, "querying invited /accounts should not return an error")

	// Unmarshal the response
	var response1 routes.AccountsCreateResponse
	err = result1.Body.FromJsonTo(&response1)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check the result's contents
	require.True(t, response1.Success, "creating a new account using inv registration failed")
	require.Equal(t, "A new account was successfully created", response1.Message, "invalid message returned by invited acc creation")
	require.NotEmpty(t, response1.Account.ID, "newly created account's id should not be empty")

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
	require.Nil(t, err, "querying invited /accounts should not return an error")

	// Unmarshal the response
	var response2 routes.AccountsCreateResponse
	err = result2.Body.FromJsonTo(&response2)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check the result's contents
	require.False(t, response2.Success, "creating a new account using invalid token should fail")
	require.Equal(t, "Invalid invitation token", response2.Message, "invalid message returned by invalid token acc creation")
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
	require.Nil(t, err, "inserting a new invitation token should not return an error")

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
	require.Nil(t, err, "querying existing invited /accounts should not return an error")

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling existing invited account creation should not return an error")

	// Check the result's contents
	require.False(t, response.Success, "creating a new account using inv registration should fail")
	require.Equal(t, "Username already exists", response.Message, "invalid message returned by existing invited acc creation")
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
	require.Nil(t, err, "inserting a new invitation token should not return an error")

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
	require.Nil(t, err, "querying expired invited /accounts should not return an error")

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling expired invited account creation should not return an error")

	// Check the result's contents
	require.False(t, response.Success, "creating a new account using inv registration should fail")
	require.Equal(t, "Expired invitation token", response.Message, "invalid message returned by expired invited acc creation")
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
	require.Nil(t, err, "inserting a new not invitation token should not return an error")

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
	require.Nil(t, err, "querying wrong type invited /accounts should not return an error")

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling wrong type invited account creation should not return an error")

	// Check the result's contents
	require.False(t, response.Success, "creating a new account using inv registration should fail")
	require.Equal(t, "Invalid invitation token", response.Message, "invalid message returned by wrong type invited acc creation")
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
	require.Nil(t, err, "querying invited /accounts should not return an error")

	// Unmarshal the response
	var createClassicResponse routes.AccountsCreateResponse
	err = createClassicResult.Body.FromJsonTo(&createClassicResponse)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check the result's contents
	require.True(t, createClassicResponse.Success, "creating a new account using classic registration failed")
	require.Equal(t, "A new account was successfully created, you should receive a confirmation email soon™.", createClassicResponse.Message, "invalid message returned by invited acc creation")
	require.NotEmpty(t, createClassicResponse.Account.ID, "newly created account's id should not be empty")
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
	require.Nil(t, err, "querying invited /accounts should not return an error")

	// Unmarshal the response
	var createClassicResponse routes.AccountsCreateResponse
	err = createClassicResult.Body.FromJsonTo(&createClassicResponse)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check the result's contents
	require.False(t, createClassicResponse.Success, "creating a new account using classic registration failed")
	require.Equal(t, "Classic registration is disabled", createClassicResponse.Message, "invalid message returned by invited acc creation")

	env.Config.ClassicRegistration = true
}

func TestAccountsCreateQueue(t *testing.T) {
	// POST /accounts - queue
	result, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			AltEmail: "something@example.com",
		},
	}.Do()
	require.Nil(t, err, "querying /accounts to queue")

	// Unmarshal the response
	var response routes.AccountsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling queue response create should not fail")

	// Check the result's contents
	require.False(t, response.Success, "creating a new account using queue registration failed")
	require.Equal(t, "Sorry, not implemented yet", response.Message, "invalid message returned by queue acc creation")
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
	require.Nil(t, err, "querying existing /tokens should not return an error")

	// Unmarshal the response
	var response routes.TokensCreateResponse
	err = request.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check the result's contents
	require.True(t, response.Success, "creating a new token using existing account failed")
	require.Equal(t, "Authentication successful", response.Message, "invalid message returned by invited acc creation")
	require.NotEmpty(t, response.Token.ID, "newly created token's id should not be empty")

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
	require.Nil(t, err, "querying /accounts should not fail")

	// Unmarshal the response
	var response routes.AccountsListResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling queue response create should not fail")

	// Check the result's contents
	require.False(t, response.Success, "creating a new account using queue registration failed")
	require.Equal(t, "Sorry, not implemented yet", response.Message, "invalid message returned by queue acc creation")
}

func TestAccountsGetMe(t *testing.T) {
	// GET /accounts/me
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err, "querying /accounts/me should not return an error")

	// Unmarshal the response
	var response routes.AccountsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling /accounts/me should not return an error")

	// Check the result's contents
	require.True(t, response.Success, "getting /accounts/me should succeed")
	require.Equal(t, "jeremy", response.Account.Name, "username should be the previously registered one")
}

func TestAccountsGetNotMe(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/not-me",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err, "querying /accounts/not-me should not return an error")

	// Unmarshal the response
	var response routes.AccountsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling /accounts/not-me should not return an error")

	// Check the result's contents
	require.False(t, response.Success, "getting /accounts/not-me should fail")
	require.Equal(t, `Only the "me" user is implemented`, response.Message, "/accounts/not-me should return a proper message")
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
	require.Nil(t, err, "updating account should not return an error")

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account update response should not return an error")

	// Check the result's contents
	require.Equal(t, "Your account has been successfully updated", response.Message, "response message should be valid")
	require.True(t, response.Success, "updating /accounts/me should succeed")
	require.Equal(t, "jeremy", response.Account.Name, "username should not be changed")
	require.Equal(t, "john.cabbage@example.com", response.Account.AltEmail, "alt email should be changed")
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
	require.Nil(t, err, "updating account should not return an error")

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account update response should not return an error")

	// Check the result's contents
	require.Equal(t, "Invalid input format", response.Message, "response message should be valid")
	require.False(t, response.Success, "updating invalid /accounts/me should fail")
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
	require.Nil(t, err, "updating account should not return an error")

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account update response should not return an error")

	// Check the result's contents
	require.Equal(t, `Only the "me" user is implemented`, response.Message, "response message should be valid")
	require.False(t, response.Success, "updating /accounts/not-me should fail")
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
	require.Nil(t, err, "updating account should not return an error")

	// Unmarshal the response
	var response routes.AccountsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account update response should not return an error")

	// Check the result's contents
	require.Equal(t, "Invalid current password", response.Message, "response message should be valid")
	require.False(t, response.Success, "updating /accounts/me should fail")
}

func TestAccountsWipeDataNotMe(t *testing.T) {
	// POST /accounts/me/wipe-data
	request := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts/not-me/wipe-data",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err, "wiping account should not return an error")

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account wipe response should not return an error")

	// Check the result's contents
	require.Equal(t, `Only the "me" user is implemented`, response.Message, "response message should be valid")
	require.False(t, response.Success, "triggering /accounts/wipe-data should fail")
}

func TestAccountsWipeData(t *testing.T) {
	// POST /accounts/me/wipe-data
	request := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts/me/wipe-data",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err, "wiping account should not return an error")

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account wipe response should not return an error")

	// Check the result's contents
	require.Equal(t, "Your account has been successfully wiped", response.Message, "response message should be valid")
	require.True(t, response.Success, "triggering /accounts/wipe-data should succeed")
}

func TestAccountsDeleteNotMe(t *testing.T) {
	// Prepare a token
	token := models.Token{
		Resource: models.MakeResource(accountID, "test invite token"),
		Type:     "auth",
	}
	token.ExpireSoon()

	err := env.Tokens.Insert(token)
	require.Nil(t, err, "inserting a new auth toekn token should not return an error")

	// DELETE /accounts/me
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/accounts/not-me",
	}
	request.AddHeader("Authorization", "Bearer "+token.ID)
	result, err := request.Do()
	require.Nil(t, err, "deleting account should not return an error")

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account delete response should not return an error")

	// Check the result's contents
	require.Equal(t, `Only the "me" user is implemented`, response.Message, "response message should be valid")
	require.False(t, response.Success, "triggering delete /account/me should fail")
}

func TestAccountsDelete(t *testing.T) {
	// Prepare a token
	token := models.Token{
		Resource: models.MakeResource(accountID, "test invite token"),
		Type:     "auth",
	}
	token.ExpireSoon()

	err := env.Tokens.Insert(token)
	require.Nil(t, err, "inserting a new auth toekn token should not return an error")

	// DELETE /accounts/me
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "Bearer "+token.ID)
	result, err := request.Do()
	require.Nil(t, err, "deleting account should not return an error")

	// Unmarshal the response
	var response routes.AccountsWipeDataResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err, "unmarshaling account delete response should not return an error")

	// Check the result's contents
	require.Equal(t, "Your account has been successfully deleted", response.Message, "response message should be valid")
	require.True(t, response.Success, "triggering delete /account/me should succeed")
}
