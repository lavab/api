package env

import (
	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"

	"github.com/lavab/api/db"
)

type Environment struct {
	Log     *logrus.Logger
	Config  *Config
	Rethink *gorethink.Session
	R       *R
}

type R struct {
	Accounts *db.AccountsTable
	Tokens   *db.TokensTable
}

var G *Environment
