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

// ContactsCreateRequest is the payload that user should pass to POST /contacts
type ContactsCreateRequest struct {
	Data         string `json:"data" schema:"data"`
	Name         string `json:"name" schema:"name"`
	Encoding     string `json:"encoding" schema:"encoding"`
	VersionMajor *int   `json:"version_major" schema:"version_major"`
	VersionMinor *int   `json:"version_minor" schema:"version_minor"`
}

// ContactsCreateResponse contains the result of the ContactsCreate request.
type ContactsCreateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Contact *models.Contact `json:"contact,omitempty"`
}

// ContactsCreate creates a new contact
func ContactsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input ContactsCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &ContactsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Fetch the current session from the database
	session := c.Env["session"].(*models.Token)

	// Ensure that the input data isn't empty
	if input.Data != "" || input.Name != "" || input.Encoding != "" || input.VersionMajor != nil || input.VersionMinor != nil {
		utils.JSONResponse(w, 400, &ContactsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Create a new contact struct
	contact := &models.Contact{
		Encrypted: models.Encrypted{
			Encoding:     input.Encoding,
			Data:         input.Data,
			Schema:       "contact",
			VersionMajor: *input.VersionMajor,
			VersionMinor: *input.VersionMinor,
		},
		Resource: models.MakeResource(session.Owner, input.Name),
	}

	// Insert the contact into the database
	if err := env.Contacts.Insert(contact); err != nil {
		utils.JSONResponse(w, 500, &ContactsCreateResponse{
			Success: false,
			Message: "internal server error - CO/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not insert a contact into the database")
		return
	}

	utils.JSONResponse(w, 201, &ContactsCreateResponse{
		Success: true,
		Message: "A new account was successfully created",
		Contact: contact,
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
