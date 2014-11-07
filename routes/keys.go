package routes

import (
	"net/http"

	"github.com/lavab/api/utils"
)

// KeysCreateResponse contains the result of the KeysCreate request.
type KeysCreateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// KeysCreate does *something* - TODO
func KeysCreate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &KeysCreateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// KeysGetResponse contains the result of the KeysGet request.
type KeysGetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// KeysGet does *something* - TODO
func KeysGet(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &KeysGetResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// KeysVoteResponse contains the result of the KeysVote request.
type KeysVoteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// KeysVote does *something* - TODO
func KeysVote(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &KeysVoteResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
