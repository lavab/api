package routes_test

import (
	"testing"

	"github.com/franela/goreq"
	"github.com/stretchr/testify/require"

	"github.com/lavab/api/routes"
)

func TestHello(t *testing.T) {
	// Request the / route
	helloResult, err := goreq.Request{
		Method: "GET",
		Uri:    server.URL,
	}.Do()
	require.Nil(t, err)

	// Unmarshal the response
	var helloResponse routes.HelloResponse
	err = helloResult.Body.FromJsonTo(&helloResponse)
	require.Nil(t, err)
	require.Equal(t, "Lavaboom API", helloResponse.Message)
}
