package routes

import (
	"fmt"
	"net/http"

	"github.com/lavab/api/utils"
)

// ContactsListResponse contains the result of the ContactsList request.
type ContactsListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactsList does *something* - TODO
func ContactsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ContactsListResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ContactsCreateResponse contains the result of the ContactsCreate request.
type ContactsCreateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactsCreate does *something* - TODO
func ContactsCreate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ContactsCreateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ContactsGetResponse contains the result of the ContactsGet request.
type ContactsGetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactsGet does *something* - TODO
func ContactsGet(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ContactsGetResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ContactsUpdateResponse contains the result of the ContactsUpdate request.
type ContactsUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactsUpdate does *something* - TODO
func ContactsUpdate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ContactsUpdateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ContactsDeleteResponse contains the result of the ContactsDelete request.
type ContactsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactsDelete does *something* - TODO
func ContactsDelete(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ContactsDeleteResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
