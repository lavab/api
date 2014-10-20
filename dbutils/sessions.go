package dbutils

import (
	"log"

	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
)

func GetSession(token string) (*models.Session, bool) {
	var result models.Session
	response, err := db.Get("sessions", token)
	if err == nil && response != nil && !response.IsNil() {
		err := response.One(&result)
		if err != nil {
			log.Fatalln("[utils.GetSession] Error when unfolding cursor")
			return nil, false
		}
		return &result, true
	}
	return nil, false
}
