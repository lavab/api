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

func TestMiddlewareNoHeader(t *testing.T) {
	result, err := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/me",
	}.Do()
	require.Nil(t, err)

	var response routes.AuthMiddlewareResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Missing auth token", response.Message)
}

func TestMiddlewareInvalidHeader(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "123")
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.AuthMiddlewareResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Invalid authorization header", response.Message)
}

func TestMiddlewareInvalidToken(t *testing.T) {
	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "Bearer 123")
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.AuthMiddlewareResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Invalid authorization token", response.Message)
}

func TestMiddlewareExpiredToken(t *testing.T) {
	// Prepare a token
	token := models.Token{
		Resource: models.MakeResource(accountID, "test invite token"),
		Expiring: models.Expiring{
			ExpiryDate: time.Now().UTC().Truncate(time.Hour * 8),
		},
		Type: "auth",
	}

	err := env.Tokens.Insert(token)
	require.Nil(t, err)

	request := goreq.Request{
		Method: "GET",
		Uri:    server.URL + "/accounts/me",
	}
	request.AddHeader("Authorization", "Bearer "+token.ID)
	result, err := request.Do()
	require.Nil(t, err)

	var response routes.AuthMiddlewareResponse
	err = result.Body.FromJsonTo(&response)
	require.Nil(t, err)

	require.False(t, response.Success)
	require.Equal(t, "Authorization token has expired", response.Message)
}
