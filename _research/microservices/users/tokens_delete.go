package main

import (
	r "github.com/dancannon/gorethink"
	"github.com/lavab/api/utils"
	"github.com/lavab/goji/web"
)

func deleteToken(c web.C, w http.ResponseWriter, req *http.Request) {
	id := c.URLParams["id"]

	wr, err := r.Table("tokens").Get(id).Delete().RunWrite(session)
	if err != nil {
		utils.JSONResponse(w, 500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.JSONResponse(w, 200, map[string]interface{}{
		"dropped": wr.Dropped,
	})
}
