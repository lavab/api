package routes

import (
	"net/http"

	"github.com/lavab/api/env"
	"github.com/lavab/api/utils"
)

type HelloResponse struct {
	Message string `json:"message"`
	DocsURL string `json:"docs_url"`
	Version string `json:"version"`
}

func Hello(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 200, &HelloResponse{
		Message: "Lavaboom API",
		DocsURL: "http://lavaboom.readme.io/",
		Version: env.G.Config.APIVersion,
	})
}
