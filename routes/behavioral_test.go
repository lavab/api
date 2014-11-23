package routes_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dancannon/gorethink"
	"github.com/franela/goreq"
	"github.com/gyepisam/mcf"
	_ "github.com/gyepisam/mcf/scrypt"
	"github.com/stretchr/testify/require"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
	"github.com/lavab/api/setup"
)

var (
	server *httptest.Server
)

func init() {
	// Mock data
	env.Config = &env.Flags{
		APIVersion:       "v0",
		LogFormatterType: "text",
		ForceColors:      false,

		SessionDuration:     72,
		ClassicRegistration: true,

		RethinkDBURL:      "127.0.0.1:28015",
		RethinkDBKey:      "",
		RethinkDBDatabase: "test",
	}

	// Connect to the RethinkDB server
	rdbSession, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:     env.Config.RethinkDBURL,
		AuthKey:     env.Config.RethinkDBKey,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})
	if err != nil {
		panic("connecting to RethinkDB should not return an error, got " + err.Error())
	}

	// Clear the test database
	err = gorethink.DbDrop("test").Exec(rdbSession)
	if err != nil {
		panic("removing the test database should not return an error, got " + err.Error())
	}

	// Disconnect
	err = rdbSession.Close()
	if err != nil {
		panic("closing the RethinkDB session should not return an error, got " + err.Error())
	}

	// Prepare a new mux (initialize the API)
	mux := setup.PrepareMux(env.Config)
	if mux == nil {
		panic("returned mux was nil")
	}

	// Set up a new temporary HTTP test server
	server = httptest.NewServer(mux)
	if server == nil {
		panic("returned httptest server was nil")
	}
}

func TestHello(t *testing.T) {
	// Request the / route
	helloResult, err := goreq.Request{
		Method: "GET",
		Uri:    server.URL,
	}.Do()
	require.Nil(t, err, "requesting / should not return an error")

	// Unmarshal the response
	var helloResponse routes.HelloResponse
	err = helloResult.Body.FromJsonTo(&helloResponse)
	require.Nil(t, err, "unmarshaling / result should not return an error")
	require.Equal(t, "Lavaboom API", helloResponse.Message)
}

func TestAccounts(t *testing.T) {
	// Constants for input data
	const (
		username = "jeremy"
		password = "potato"
	)

	// Hash the passwords
	passwordHashed, err := mcf.Create(password)
	require.Nil(t, err, "hashing password should not return an error")

	// Prepare a token
	inviteToken := models.Token{
		Resource: models.MakeResource("", "test invite token"),
		Type:     "invite",
	}
	inviteToken.ExpireSoon()
	err = env.Tokens.Insert(inviteToken)
	require.Nil(t, err, "inserting a new invitation token should not return an error")

	// POST /accounts - unknown
	createUnknownResult, err := goreq.Request{
		Method: "POST",
		Uri:    server.URL + "/accounts",
	}.Do()
	require.Nil(t, err, "querying unknown /accounts should not return an error")

	// Unmarshal the response
	var createUnknownResponse routes.AccountsCreateResponse
	err = createUnknownResult.Body.FromJsonTo(&createUnknownResponse)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check values
	require.False(t, createUnknownResponse.Success, "unknown request should return success false")
	require.Equal(t, "Invalid request", createUnknownResponse.Message, "unknown request should return proper error msg")

	// POST /accounts - invited
	createInvitedResult, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username + "invited",
			Password: passwordHashed,
			Token:    inviteToken.ID,
		},
	}.Do()
	require.Nil(t, err, "querying invited /accounts should not return an error")

	// Unmarshal the response
	var createInvitedResponse routes.AccountsCreateResponse
	err = createInvitedResult.Body.FromJsonTo(&createInvitedResponse)
	require.Nil(t, err, "unmarshaling invited account creation should not return an error")

	// Check the result's contents
	require.True(t, createInvitedResponse.Success, "creating a new account using inv registration failed")
	require.Equal(t, "A new account was successfully created", createInvitedResponse.Message, "invalid message returned by invited acc creation")
	require.NotEmpty(t, createInvitedResponse.Account.ID, "newly created account's id should not be empty")

	// POST /accounts - classic
	createClassicResult, err := goreq.Request{
		Method:      "POST",
		Uri:         server.URL + "/accounts",
		ContentType: "application/json",
		Body: routes.AccountsCreateRequest{
			Username: username + "classic",
			Password: passwordHashed,
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
