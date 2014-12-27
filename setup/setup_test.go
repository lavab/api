package setup

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lavab/api/env"
)

func TestSetup(t *testing.T) {
	// Mock data
	env.Config = &env.Flags{
		APIVersion:       "v0",
		LogFormatterType: "text",
		ForceColors:      true,

		SessionDuration:     72,
		ClassicRegistration: true,

		RedisAddress: "127.0.0.1:6379",

		NATSAddress: "nats://127.0.0.1:4222",

		RethinkDBAddress:  "127.0.0.1:28015",
		RethinkDBKey:      "",
		RethinkDBDatabase: "test",
	}

	// Prepare a new mux (initialize the API)
	mux := PrepareMux(env.Config)
	require.NotNil(t, mux, "mux should not be nil")
}
