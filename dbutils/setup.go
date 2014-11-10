package dbutils

import (
	"github.com/lavab/api/db"
)

//The crud interfaces for models
var Users UsersTable
var Sessions SessionTable

//The init version for routes, mostly related to
//initing the table informatio
func init() {
	//at this stage we should have the db config variable initialized
	userCrud := db.NewCrudTable(db.CurrentConfig.Db, db.TABLE_USERS)
	Users = UsersTable{RethinkCrud: userCrud}

	//init the sessions variable
	sessionCrud := db.NewCrudTable(db.CurrentConfig.Db, db.TABLE_SESSIONS)
	Sessions = SessionTable{RethinkCrud: sessionCrud}

}
