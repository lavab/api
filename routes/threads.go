package routes

import (
	"net/http"

	"github.com/lavab/api/utils"
)

// ThreadsListResponse contains the result of the ThreadsList request.
type ThreadsListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ThreadsList shows all threads
func ThreadsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ThreadsListResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ThreadsGetResponse contains the result of the ThreadsGet request.
type ThreadsGetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ThreadsGet returns information about a single thread.
func ThreadsGet(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ThreadsGetResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ThreadsUpdateResponse contains the result of the ThreadsUpdate request.
type ThreadsUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ThreadsUpdate does *something* with a thread.
func ThreadsUpdate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ThreadsUpdateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
