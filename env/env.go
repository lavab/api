package env

import (
	"github.com/Sirupsen/logrus"
)

type Environment struct {
	Log    *logrus.Logger
	Config *Config
}

var G *Environment
