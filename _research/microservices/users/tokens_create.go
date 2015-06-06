package main

import (
	"encoding/json"

	r "github.com/dancannon/gorethink"
	"github.com/dchest/uniuri"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
	"github.com/lavab/goji/web"
)

func deleteToken(w http.ResponseWriter, req *http.Request) {
	var user *models.Token
	if err := json.NewDecoder(req).Decode(&user); err != nil {
		utils.JSONResponse(w, 500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	user.ID = uniuri.NewLen(uniuri.UUIDLen)

	wr, err := r.Table("tokens").Insert(user).RunWrite(session)
	if err != nil {
		utils.JSONResponse(w, 500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.JSONResponse(w, 200, map[string]interface{}{
		"created": wr.Created,
	})
}
