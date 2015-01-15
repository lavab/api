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
	emailID         string
	notOwnedEmailID string
)

func TestEmailsCreate(t *testing.T) {
	request := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/emails",
		ContentType: "application/json",
		Body: routes.EmailsCreateRequest{
			To:      []string{"piotr@zduniak.net"},
			Subject: "hello world",
			Body:    "raw meaty email",
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, len(response.Created), 1)
	require.True(t, response.Success)

	emailID = response.Created[0]
}

func TestEmailsCreateInvalidBody(t *testing.T) {
	request := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/emails",
		ContentType: "application/json",
		Body:        "!@#!@#!@#",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, response.Message, "Invalid input format")
	require.False(t, response.Success)
}

func TestEmailsCreateMissingFields(t *testing.T) {
	request := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/emails",
		ContentType: "application/json",
		Body: routes.EmailsCreateRequest{
			To:      []string{"piotr@zduniak.net"},
			Subject: "hello world",
		},
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsCreateResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, response.Message, "Invalid request")
	require.False(t, response.Success)
}

func TestEmailsGet(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/emails/" + emailID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.Equal(t, "hello world", response.Email.Name)
}

func TestEmailsGetInvalidID(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/emails/nonexisting",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, "Email not found", response.Message)
	require.False(t, response.Success)
}

func TestEmailsGetNotOwned(t *testing.T) {
	email := &models.Email{
		Resource: models.MakeResource("not", "Carpeus Caesar"),
	}

	err := env.Emails.Insert(email)
	require.Nil(t, err)

	notOwnedEmailID = email.ID

	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/emails/" + email.ID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsGetResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.Equal(t, "Email not found", response.Message)
	require.False(t, response.Success)
}

func TestEmailsList(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/emails?offset=0&limit=1&sort=+date",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsListResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.Equal(t, "hello world", (*response.Emails)[0].Name)
}

func TestEmailsListInvalidOffset(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/emails?offset=pppp&limit=1&sort=+date",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsListResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Invalid offset", response.Message)
}

func TestEmailsListInvalidLimit(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/emails?offset=0&limit=pppp&sort=+date",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsListResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Invalid limit", response.Message)
}

func TestEmailsDelete(t *testing.T) {
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/emails/" + emailID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsDeleteResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.True(t, response.Success)
	require.Equal(t, "Email successfully removed", response.Message)
}

func TestEmailsDeleteNotExisting(t *testing.T) {
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/emails/nonexisting",
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsDeleteResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Email not found", response.Message)
}

func TestEmailsDeleteNotOwned(t *testing.T) {
	request := goreq.Request{
		Method: "DELETE",
		Uri:    server.URL + "/emails/" + notOwnedEmailID,
	}
	request.AddHeader("Authorization", "Bearer "+authToken)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.EmailsDeleteResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Email not found", response.Message)
}
