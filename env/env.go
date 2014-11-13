package env

import (
	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/db"
)

var (
	// Config contains flags passed to the API
	Config *Flags
	// Log is the API's logrus instance
	Log *logrus.Logger
	// Rethink contains the RethinkDB session used in the API
	Rethink *gorethink.Session
	// Accounts is the global instance of AccountsTable
	Accounts *db.AccountsTable
	// Tokens is the global instance of TokensTable
	Tokens *db.TokensTable
	// Keys is the global instance of KeysTable
	Keys *db.KeysTable
	// TokenCache is an instance of TokenCache interface to
	// keep the tokens in Redis instance
	TokensCache cache.TokenCache
)
