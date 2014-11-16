package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// ContactsListResponse contains the result of the ContactsList request.
type ContactsListResponse struct {
	Success  bool               `json:"success"`
	Message  string             `json:"message,omitempty"`
	Contacts *[]*models.Contact `json:"contacts,omitempty"`
}

// ContactsList does *something* - TODO
func ContactsList(c web.C, w http.ResponseWriter, r *http.Request) {
	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Get contacts from the database
	contacts, err := env.Contacts.GetOwnedBy(session.Owner)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to fetch contacts")

		utils.JSONResponse(w, 500, &AccountsDeleteResponse{
			Success: false,
			Message: "Internal error (code CO/LI/01)",
		})
		return
	}

	utils.JSONResponse(w, 501, &ContactsListResponse{
		Success:  false,
		Contacts: &contacts,
	})
}

type ContactsCreateRequest struct {
	Data         string `json:"data" schema:"data"`
	Name         string `json:"name" schema:"name"`
	Encoding     string `json:"encoding" schema:"encoding"`
	VersionMajor int    `json:"version_major" schema:"version_major"`
	VersionMinor int    `json:"version_minor" schema:"version_minor"`
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
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Contact *models.Contact `json:"contact"`
}

// ContactsGet does *something* - TODO
func ContactsGet(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ContactsGetResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

type ContactsUpdateRequest struct {
	Data         string `json:"data" schema:"data"`
	Name         string `json:"name" schema:"name"`
	Encoding     string `json:"encoding" schema:"encoding"`
	VersionMajor int    `json:"version_major" schema:"version_major"`
	VersionMinor int    `json:"version_minor" schema:"version_minor"`
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
