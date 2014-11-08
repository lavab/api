package dbutils

import (
	"log"

	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
)

type SessionTable struct {
	db.RethinkCrud
}

func (sessions *SessionTable) GetSession(id string) (*models.Session, bool) {
	var result models.Session

	if err := sessions.FindFetchOne(id, &result); err != nil {
		log.Println(err.Error())
		return nil, false
	}

	return &result, true
}
