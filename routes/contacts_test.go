package routes_test

import (
	"testing"

	"github.com/franela/goreq"
	"github.com/stretchr/testify/require"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
)

func TestContactsPrepareAccount(t *testing.T) {
	const (
		username = "jeremy-contacts"
		password = "potato"
	)

	inviteToken := models.Token{
		Resource: models.MakeResource("", "test invite token"),
		Type:     "invite",
	}
	inviteToken.ExpireSoon()

	err := env.Tokens.Insert(inviteToken)
	require.Nil(t, err)

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

	var response1 routes.AccountsCreateResponse
	err = result1.Body.FromJsonTo(&response1)
	require.Nil(t, err)

	require.True(t, response1.Success)
	require.Equal(t, "A new account was successfully created", response1.Message)
	require.NotEmpty(t, response1.Account.ID)

	accountID = response1.Account.ID

	request2, err := goreq.Request{
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

	var response2 routes.TokensCreateResponse
	err = request2.Body.FromJsonTo(&response2)
	require.Nil(t, err)

	require.True(t, response2.Success)
	require.Equal(t, "Authentication successful", response2.Message)
	require.NotEmpty(t, response2.Token.ID)

	authToken = response2.Token.ID
}

func TestContactsCreate(t *testing.T) {
	request := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/contacts",
		ContentType: "application/json",
		Body: routes.ContactsCreateRequest{
			Data:         "random stuff",
			Name:         "John Doe",
			Encoding:     "json",
			VersionMajor: 1,
			VersionMinor: 0,
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, "A new contact was successfully created", response.Message)
	require.True(t, response.Success)
	require.NotEmpty(t, response.Contact.ID)
}

func TestContactsCreateMissingParts(t *testing.T) {
	request := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/contacts",
		ContentType: "application/json",
		Body: routes.ContactsCreateRequest{
			Data:         "random stuff",
			Name:         "John Doe",
			Encoding:     "",
			VersionMajor: 1,
			VersionMinor: 0,
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, "Invalid request", response.Message)
	require.False(t, response.Success)
}

func TestContactsCreateInvalid(t *testing.T) {
	request := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/contacts",
		ContentType: "application/json",
		Body:        "!@#!@#!@#",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, "Invalid input format", response.Message)
	require.False(t, response.Success)
}

func TestContactsList(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/contacts",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsListResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.True(t, len(*response.Contacts) > 0)
}
