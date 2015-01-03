package routes_test

import (
	"testing"

	"github.com/franela/goreq"
	"github.com/stretchr/testify/require"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
)

var (
	contactID         string
	notOwnedContactID string
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
			Data:            "random stuff",
			Name:            "John Doe",
			Encoding:        "json",
			VersionMajor:    1,
			VersionMinor:    0,
			PGPFingerprints: []string{"that's totally a fingerprint!"},
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

	contactID = response.Contact.ID
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

func TestContactsGet(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/contacts/" + contactID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.Equal(t, "John Doe", response.Contact.Name)
}

func TestContactsGetNotOwned(t *testing.T) {
	contact := &models.Contact{
		Encrypted: models.Encrypted{
			Encoding:     "json",
			Data:         "carp",
			Schema:       "contact",
			VersionMajor: 1,
			VersionMinor: 0,
		},
		Resource: models.MakeResource("not", "Carpeus Caesar"),
	}

	err := env.Contacts.Insert(contact)
	require.Nil(t, err)

	notOwnedContactID = contact.ID

	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/contacts/" + contact.ID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Contact not found", response.Message)
}

func TestContactsGetWrongID(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/contacts/gibberish",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Contact not found", response.Message)
}

func TestContactsUpdate(t *testing.T) {
	int1 := 1
	int0 := 0

	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/contacts/" + contactID,
		ContentType: "application/json",
		Body: routes.ContactsUpdateRequest{
			Data:         "random stuff2",
			Name:         "John Doez",
			Encoding:     "json",
			VersionMajor: &int1,
			VersionMinor: &int0,
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.Equal(t, "John Doez", response.Contact.Name)
}

func TestContactsUpdateInvalid(t *testing.T) {
	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/contacts/" + contactID,
		ContentType: "application/json",
		Body:        "123123!@#!@#",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Invalid input format", response.Message)
}

func TestContactsUpdateNotOwned(t *testing.T) {
	int1 := 1
	int0 := 0

	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/contacts/" + notOwnedContactID,
		ContentType: "application/json",
		Body: routes.ContactsUpdateRequest{
			Data:         "random stuff2",
			Name:         "John Doez",
			Encoding:     "json",
			VersionMajor: &int1,
			VersionMinor: &int0,
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Contact not found", response.Message)
}

func TestContactsUpdateNotExisting(t *testing.T) {
	int1 := 1
	int0 := 0

	request := goreq.Request{
		Method:      "PUT",
		Uri:         server.URL + "/contacts/gibberish",
		ContentType: "application/json",
		Body: routes.ContactsUpdateRequest{
			Data:         "random stuff2",
			Name:         "John Doez",
			Encoding:     "json",
			VersionMajor: &int1,
			VersionMinor: &int0,
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsUpdateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Contact not found", response.Message)
}

func TestContactsDelete(t *testing.T) {
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/contacts/" + contactID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsDeleteResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.Equal(t, "Contact successfully removed", response.Message)
}

func TestContactsDeleteNotOwned(t *testing.T) {
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/contacts/" + notOwnedContactID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsDeleteResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Contact not found", response.Message)
}

func TestContactsDeleteNotExisting(t *testing.T) {
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/contacts/gibberish",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.ContactsDeleteResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Contact not found", response.Message)
}
