package dbutils

import (
	"log"

	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
)

// TODO change names to auth tokens instead of sessions
func GetSession(id string) (*models.AuthToken, bool) {
	var result models.AuthToken
	response, err := db.Get("sessions", id)
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
