package main

import (
	"encoding/json"

	r "github.com/dancannon/gorethink"
	"github.com/dchest/uniuri"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
	"github.com/lavab/goji/web"
)

func updateToken(c web.C, w http.ResponseWriter, req *http.Request) {
	var token *models.Token
	if err := json.NewDecoder(req).Decode(&token); err != nil {
		utils.JSONResponse(w, 400, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	token.ID = uniuri.NewLen(uniuri.UUIDLen)

	wr, err := r.Table("tokens").Insert(token).RunWrite(session)
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
