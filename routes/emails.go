package routes

import (
	"net/http"

	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// EmailsListResponse contains the result of the EmailsList request.
type EmailsListResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message,omitempty"`
	ItemsCount int             `json:"items_count,omitempty"`
	Emails     []*models.Email `json:"emails,omitempty"`
}

// EmailsList sends a list of the emails in the inbox.
func EmailsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 200, &EmailsListResponse{
		Success:    true,
		ItemsCount: 1,
		Emails:     []*models.Email{},
	})
}

// EmailsCreateResponse contains the result of the EmailsCreate request.
type EmailsCreateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Created []string `json:"created,omitempty"`
}

// EmailsCreate sends a new email
func EmailsCreate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 200, &EmailsCreateResponse{
		Success: true,
		Created: []string{"123"},
	})
}

// EmailsGetResponse contains the result of the EmailsGet request.
type EmailsGetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

// EmailsGet responds with a single email message
func EmailsGet(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 200, &EmailsGetResponse{
		Success: true,
		Status:  "sending",
	})
}

// EmailsUpdateResponse contains the result of the EmailsUpdate request.
type EmailsUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EmailsUpdate does *something* - TODO
func EmailsUpdate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &EmailsUpdateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// EmailsDeleteResponse contains the result of the EmailsDelete request.
type EmailsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EmailsDelete remvoes an email from the system
func EmailsDelete(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &EmailsDeleteResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
