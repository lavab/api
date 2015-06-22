package routes

import (
	"net/http"
	"strings"

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
		query         = req.URL.Query()
		sTags         = query.Get("tags")
		excludeBodies = query.Get("exclude_bodies")
		result        []*models.File
	)

	q := r.Table("files")

	if sTags == "" {
		q = q.GetAllByIndex("owner", session.Owner)
	} else {
		tags := strings.Split(sTags, ",")
		ids := []interface{}{}
		for _, tag := range tags {
			ids = append(ids, []interface{}{
				session.Owner,
				tag,
			})
		}

		q = q.GetAllByIndex("ownerTags", ids...)
	}

	if excludeBodies == "true" {
		q = q.Without("body")
	}

	cursor, err := q.Run(env.Rethink)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.FilesListUnableToGet, err, false,
		))
		return
	}
	defer cursor.Close()
	if err := cursor.All(&result); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.FilesListUnableToGet, err, false,
		))
		return
	}

	utils.JSONResponse(w, 200, &FilesListResponse{
		Success: true,
		Files:   &result,
	})
}

type FilesCreateRequest struct {
	Name string                 `json:"name" schema:"name"`
	Meta map[string]interface{} `json:"meta" schema:"meta"`
	Body []byte                 `json:"body" schema:"body"`
	Tags []string               `json:"tags" schema:"tags"`
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.FilesCreateInvalidInput, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Create a new file struct
	file := &models.File{
		Resource: models.MakeResource(session.Owner, input.Name),
		Meta:     input.Meta,
		Body:     input.Body,
		Tags:     input.Tags,
	}

	// Insert the file into the database
	if err := env.Files.Insert(file); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.FilesCreateUnableToInsert, err, true,
		))
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
	cursor, err := r.Table("files").Get(c.URLParams["id"]).Run(env.Rethink)
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.FilesGetUnableToGet, err, true,
		))
		return
	}
	defer cursor.Close()
	var file *models.File
	if err := cursor.One(&file); err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.FilesGetUnableToGet, err, true,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if file.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.FilesGetNotOwned, "You're not the owner of this file", false,
		))
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
	Name *string                `json:"name" schema:"name"`
	Meta map[string]interface{} `json:"meta" schema:"meta"`
	Body []byte                 `json:"body" schema:"body"`
	Tags []string               `json:"tags" schema:"tags"`
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.FilesUpdateInvalidInput, err, false,
		))
		return
	}

	// Get the file from the database
	file, err := env.Files.GetFile(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.FilesUpdateUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if file.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.FilesUpdateNotOwned, "You're not the owner of this file", false,
		))
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
		utils.JSONResponse(w, 500, utils.NewError(
			utils.FilesUpdateUnableToUpdate, err, true,
		))
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
		utils.JSONResponse(w, 404, utils.NewError(
			utils.FilesDeleteUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if file.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.FilesDeleteNotOwned, "You're not the owner of this file", false,
		))
		return
	}

	// Perform the deletion
	err = env.Files.DeleteID(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.FilesDeleteUnableToDelete, err, true,
		))
		return
	}

	// Write the file to the response
	utils.JSONResponse(w, 200, &FilesDeleteResponse{
		Success: true,
		Message: "File successfully removed",
	})
}
