package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

type FilesListResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Files   *[]*models.File `json:"files,omitempty"`
}

func FilesList(c web.C, w http.ResponseWriter, r *http.Request) {
	session := c.Env["token"].(*models.Token)

	query := r.URL.Query()
	email := query.Get("email")
	name := query.Get("name")

	if email == "" || name == "" {
		utils.JSONResponse(w, 400, &FilesListResponse{
			Success: false,
			Message: "No email or name in get params",
		})
		return
	}

	files, err := env.Files.GetInEmail(session.Owner, email, name)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch files")

		utils.JSONResponse(w, 500, &FilesListResponse{
			Success: false,
			Message: "Internal error (code FI/LI/01)",
		})
		return
	}

	utils.JSONResponse(w, 200, &FilesListResponse{
		Success: true,
		Files:   &files,
	})
}

type FilesCreateRequest struct {
	Data            string   `json:"data" schema:"data"`
	Name            string   `json:"name" schema:"name"`
	Encoding        string   `json:"encoding" schema:"encoding"`
	VersionMajor    int      `json:"version_major" schema:"version_major"`
	VersionMinor    int      `json:"version_minor" schema:"version_minor"`
	PGPFingerprints []string `json:"pgp_fingerprints" schema:"pgp_fingerprints"`
}

type FilesCreateResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	File    *models.File `json:"file,omitempty"`
}

// FilesCreate creates a new file
func FilesCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input FilesCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &FilesCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the input data isn't empty
	if input.Data == "" || input.Name == "" || input.Encoding == "" {
		utils.JSONResponse(w, 400, &FilesCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Create a new file struct
	file := &models.File{
		Encrypted: models.Encrypted{
			Encoding:        input.Encoding,
			Data:            input.Data,
			Schema:          "file",
			VersionMajor:    input.VersionMajor,
			VersionMinor:    input.VersionMinor,
			PGPFingerprints: input.PGPFingerprints,
		},
		Resource: models.MakeResource(session.Owner, input.Name),
	}

	// Insert the file into the database
	if err := env.Files.Insert(file); err != nil {
		utils.JSONResponse(w, 500, &FilesCreateResponse{
			Success: false,
			Message: "internal server error - FI/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not insert a file into the database")
		return
	}

	utils.JSONResponse(w, 201, &FilesCreateResponse{
		Success: true,
		Message: "A new file was successfully created",
		File:    file,
	})
}

// FilesGetResponse contains the result of the FilesGet request.
type FilesGetResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	File    *models.File `json:"file,omitempty"`
}

// FilesGet gets the requested file from the database
func FilesGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the file from the database
	file, err := env.Files.GetFile(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &FilesGetResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if file.Owner != session.Owner {
		utils.JSONResponse(w, 404, &FilesGetResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Write the file to the response
	utils.JSONResponse(w, 200, &FilesGetResponse{
		Success: true,
		File:    file,
	})
}

// FilesUpdateRequest is the payload passed to PUT /files/:id
type FilesUpdateRequest struct {
	Data            string   `json:"data" schema:"data"`
	Name            string   `json:"name" schema:"name"`
	Encoding        string   `json:"encoding" schema:"encoding"`
	VersionMajor    *int     `json:"version_major" schema:"version_major"`
	VersionMinor    *int     `json:"version_minor" schema:"version_minor"`
	PGPFingerprints []string `json:"pgp_fingerprints" schema:"pgp_fingerprints"`
}

// FilesUpdateResponse contains the result of the FilesUpdate request.
type FilesUpdateResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	File    *models.File `json:"file,omitempty"`
}

// FilesUpdate updates an existing file in the database
func FilesUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input FilesUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &FilesUpdateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the file from the database
	file, err := env.Files.GetFile(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &FilesUpdateResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if file.Owner != session.Owner {
		utils.JSONResponse(w, 404, &FilesUpdateResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	if input.Data != "" {
		file.Data = input.Data
	}

	if input.Name != "" {
		file.Name = input.Name
	}

	if input.Encoding != "" {
		file.Encoding = input.Encoding
	}

	if input.VersionMajor != nil {
		file.VersionMajor = *input.VersionMajor
	}

	if input.VersionMinor != nil {
		file.VersionMinor = *input.VersionMinor
	}

	if input.PGPFingerprints != nil {
		file.PGPFingerprints = input.PGPFingerprints
	}

	// Perform the update
	err = env.Files.UpdateID(c.URLParams["id"], file)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to update a file")

		utils.JSONResponse(w, 500, &FilesUpdateResponse{
			Success: false,
			Message: "Internal error (code FI/UP/01)",
		})
		return
	}

	// Write the file to the response
	utils.JSONResponse(w, 200, &FilesUpdateResponse{
		Success: true,
		File:    file,
	})
}

// FilesDeleteResponse contains the result of the Delete request.
type FilesDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// FilesDelete removes a file from the database
func FilesDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the file from the database
	file, err := env.Files.GetFile(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &FilesDeleteResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if file.Owner != session.Owner {
		utils.JSONResponse(w, 404, &FilesDeleteResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Perform the deletion
	err = env.Files.DeleteID(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete a file")

		utils.JSONResponse(w, 500, &FilesDeleteResponse{
			Success: false,
			Message: "Internal error (code FI/DE/01)",
		})
		return
	}

	// Write the file to the response
	utils.JSONResponse(w, 200, &FilesDeleteResponse{
		Success: true,
		Message: "File successfully removed",
	})
}
