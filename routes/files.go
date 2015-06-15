package routes

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
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

func FilesList(c web.C, w http.ResponseWriter, req *http.Request) {
	session := c.Env["token"].(*models.Token)

	var (
		query  = req.URL.Query()
		sTags  = query.Get("tags")
		result []*models.File
	)

	if sTags == "" {
		cursor, err := r.Table("files").GetAllByIndex("owner", session.Owner).Run(env.Rethink)
		if err != nil {
			utils.JSONResponse(w, 500, &FilesListResponse{
				Success: false,
				Message: "Internal error (code FI/LI/01)",
			})
			return
		}
		defer cursor.Close()
		if err := cursor.All(&result); err != nil {
			utils.JSONResponse(w, 500, &FilesListResponse{
				Success: false,
				Message: "Internal error (code FI/LI/02)",
			})
			return
		}
	} else {
		tags := strings.Split(sTags, ",")
		ids := []interface{}{}
		for _, tag := range tags {
			ids = append(ids, []interface{}{
				session.Owner,
				tag,
			})
		}
		cursor, err := r.Table("files").GetAllByIndex("ownerTags", ids...).Run(env.Rethink)
		if err != nil {
			utils.JSONResponse(w, 500, &FilesListResponse{
				Success: false,
				Message: "Internal error (code FI/LI/03)",
			})
			return
		}
		defer cursor.Close()
		if err := cursor.All(&result); err != nil {
			utils.JSONResponse(w, 500, &FilesListResponse{
				Success: false,
				Message: "Internal error (code FI/LI/04)",
			})
			return
		}
	}

	utils.JSONResponse(w, 200, &FilesListResponse{
		Success: true,
		Files:   &result,
	})
}

type FilesCreateRequest struct {
	Name string      `json:"name" schema:"name"`
	Meta interface{} `json:"meta" schema:"meta"`
	Body string      `json:"body" schema:"body"`
	Tags []string    `json:"tags" schema:"tags"`
}

type FilesCreateResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	ID      *string `json:"file,omitempty"`
}

// FilesCreate creates a new file
func FilesCreate(c web.C, w http.ResponseWriter, req *http.Request) {
	// Decode the request
	var input FilesCreateRequest
	err := utils.ParseRequest(req, &input)
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

	// Decode the body
	body, err := base64.StdEncoding.DecodeString(input.Body)
	if err != nil {
		utils.JSONResponse(w, 400, &FilesCreateResponse{
			Success: false,
			Message: "Invalid input format, " + err.Error(),
		})
		return
	}

	// Create a new file struct
	file := &models.File{
		Resource: models.MakeResource(session.Owner, input.Name),
		Meta:     input.Meta,
		Body:     body,
		Tags:     input.Tags,
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
		ID:      &file.ID,
	})
}

// FilesGetResponse contains the result of the FilesGet request.
type FilesGetResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	File    *models.File `json:"file,omitempty"`
}

// FilesGet gets the requested file from the database
func FilesGet(c web.C, w http.ResponseWriter, req *http.Request) {
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
	Name *string     `json:"name" schema:"name"`
	Meta interface{} `json:"meta" schema:"meta"`
	Body []byte      `json:"body" schema:"body"`
	Tags []string    `json:"tags" schema:"tags"`
}

// FilesUpdateResponse contains the result of the FilesUpdate request.
type FilesUpdateResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	File    *models.File `json:"file,omitempty"`
}

// FilesUpdate updates an existing file in the database
func FilesUpdate(c web.C, w http.ResponseWriter, req *http.Request) {
	// Decode the request
	var input FilesUpdateRequest
	err := utils.ParseRequest(req, &input)
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

	if input.Name != nil {
		file.Name = *input.Name
	}

	if input.Meta != nil {
		file.Meta = input.Meta
	}

	if input.Body != nil {
		file.Body = input.Body
	}

	if input.Tags != nil {
		file.Tags = input.Tags
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
func FilesDelete(c web.C, w http.ResponseWriter, req *http.Request) {
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
