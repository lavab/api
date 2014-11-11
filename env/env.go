package env

import (
	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"

	"github.com/lavab/api/db"
)

var (
	Config   *Flags
	Log      *logrus.Logger
	Rethink  *gorethink.Session
	Accounts *db.AccountsTable
	Tokens   *db.TokensTable
)
