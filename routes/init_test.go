package routes_test

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/dancannon/gorethink"

	"github.com/lavab/api/env"
	"github.com/lavab/api/setup"
)

var (
	server    *httptest.Server
	authToken string
)

func init() {
	// Mock data
	env.Config = &env.Flags{
		APIVersion:       "v0",
		LogFormatterType: "text",
		ForceColors:      true,

		SessionDuration: 72,

		RedisAddress: "127.0.0.1:6379",

		NSQdAddress:    "127.0.0.1:4150",
		LookupdAddress: "127.0.0.1:4160",

		RethinkDBAddress:  "127.0.0.1:28015",
		RethinkDBKey:      "",
		RethinkDBDatabase: "test",
	}

	// Connect to the RethinkDB server
	rdbSession, err := gorethink.Connect(gorethink.ConnectOpts{
		Address: env.Config.RethinkDBAddress,
		AuthKey: env.Config.RethinkDBKey,
		MaxIdle: 10,
		Timeout: time.Second * 10,
	})
	if err != nil {
		panic("connecting to RethinkDB should not return an error, got " + err.Error())
	}

	// Clear the test database
	err = gorethink.DbDrop("test").Exec(rdbSession)
	if err != nil {
		fmt.Println("removing the test database should not return an error, got " + err.Error())
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
