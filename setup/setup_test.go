package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lavab/api/env"
)

func TestSetup(t *testing.T) {
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

	// Prepare a new mux (initialize the API)
	mux := PrepareMux(env.Config)
	assert.NotNil(t, mux, "mux should not be nil")
}
