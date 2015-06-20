package routes

import (
	"net/http"

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
	session := c.Env["token"].(*models.Token)

	// Get contacts from the database
	contacts, err := env.Contacts.GetOwnedBy(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ContactsListUnableToGet, err, true,
		))
		return
	}

	utils.JSONResponse(w, 200, &ContactsListResponse{
		Success:  true,
		Contacts: &contacts,
	})
}

// ContactsCreateRequest is the payload that user should pass to POST /contacts
type ContactsCreateRequest struct {
	Data         string `json:"data" schema:"data"`
	Name         string `json:"name" schema:"name"`
	Encoding     string `json:"encoding" schema:"encoding"`
	VersionMajor int    `json:"version_major" schema:"version_major"`
	VersionMinor int    `json:"version_minor" schema:"version_minor"`
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.ContactsCreateInvalidInput, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the input data isn't empty
	if input.Data == "" || input.Name == "" || input.Encoding == "" {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.ContactsCreateInvalidInput, "One of the input fields is empty", false,
		))
		return
	}

	// Create a new contact struct
	contact := &models.Contact{
		Encrypted: models.Encrypted{
			Encoding:     input.Encoding,
			Data:         input.Data,
			Schema:       "contact",
			VersionMajor: input.VersionMajor,
			VersionMinor: input.VersionMinor,
		},
		Resource: models.MakeResource(session.Owner, input.Name),
	}

	// Insert the contact into the database
	if err := env.Contacts.Insert(contact); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ContactsCreateUnableToInsert, err, true,
		))
		return
	}

	utils.JSONResponse(w, 201, &ContactsCreateResponse{
		Success: true,
		Message: "A new contact was successfully created",
		Contact: contact,
	})
}

// ContactsGetResponse contains the result of the ContactsGet request.
type ContactsGetResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Contact *models.Contact `json:"contact,omitempty"`
}

// ContactsGet gets the requested contact from the database
func ContactsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the contact from the database
	contact, err := env.Contacts.GetContact(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ContactsGetUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if contact.Owner != session.Owner {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.ContactsGetNotOwned, "You're not the owner of this contact", false,
		))
		return
	}

	// Write the contact to the response
	utils.JSONResponse(w, 200, &ContactsGetResponse{
		Success: true,
		Contact: contact,
	})
}

// ContactsUpdateRequest is the payload passed to PUT /contacts/:id
type ContactsUpdateRequest struct {
	Data         string `json:"data" schema:"data"`
	Name         string `json:"name" schema:"name"`
	Encoding     string `json:"encoding" schema:"encoding"`
	VersionMajor *int   `json:"version_major" schema:"version_major"`
	VersionMinor *int   `json:"version_minor" schema:"version_minor"`
}

// ContactsUpdateResponse contains the result of the ContactsUpdate request.
type ContactsUpdateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Contact *models.Contact `json:"contact,omitempty"`
}

// ContactsUpdate updates an existing contact in the database
func ContactsUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input ContactsUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.ContactsUpdateInvalidInput, err, false,
		))
		return
	}

	// Get the contact from the database
	contact, err := env.Contacts.GetContact(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ContactsUpdateUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if contact.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ContactsUpdateNotOwned, "You're not the owner of this contact", false,
		))
		return
	}

	if input.Data != "" {
		contact.Data = input.Data
	}

	if input.Name != "" {
		contact.Name = input.Name
	}

	if input.Encoding != "" {
		contact.Encoding = input.Encoding
	}

	if input.VersionMajor != nil {
		contact.VersionMajor = *input.VersionMajor
	}

	if input.VersionMinor != nil {
		contact.VersionMinor = *input.VersionMinor
	}

	// Perform the update
	err = env.Contacts.UpdateID(c.URLParams["id"], contact)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ContactsUpdateUnableToUpdate, err, true,
		))
		return
	}

	// Write the contact to the response
	utils.JSONResponse(w, 200, &ContactsUpdateResponse{
		Success: true,
		Contact: contact,
	})
}

// ContactsDeleteResponse contains the result of the ContactsDelete request.
type ContactsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactsDelete removes a contact from the database
func ContactsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the contact from the database
	contact, err := env.Contacts.GetContact(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ContactsDeleteUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if contact.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ContactsDeleteNotOwned, "You're not the owner of this contact", false,
		))
		return
	}

	// Perform the deletion
	err = env.Contacts.DeleteID(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ContactsDeleteUnableToDelete, "You're not the owner of this contact", true,
		))
		return
	}

	// Write the contact to the response
	utils.JSONResponse(w, 200, &ContactsDeleteResponse{
		Success: true,
		Message: "Contact successfully removed",
	})
}
