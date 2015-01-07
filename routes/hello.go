package routes

import (
	"net/http"

	"github.com/lavab/api/env"
	"github.com/lavab/api/utils"
)

// HelloResponse contains the result of the Hello request.
type HelloResponse struct {
	Message string `json:"message"`
	DocsURL string `json:"docs_url"`
	Version string `json:"version"`
}

// Hello shows basic information about the API on its frontpage.
func Hello(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 200, &HelloResponse{
		Message: "Lavaboom API",
		DocsURL: "https://docs.lavaboom.io/",
		Version: env.Config.APIVersion,
	})
}
