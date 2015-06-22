package routes

import (
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
	"github.com/zenazn/goji/web"
)

// LabelsListResponse contains the result of the LabelsList request.
type LabelsListResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Labels  *[]*models.Label `json:"labels,omitempty"`
}

// LabelsList fetches all labels
func LabelsList(c web.C, w http.ResponseWriter, req *http.Request) {
	session := c.Env["token"].(*models.Token)

	cursor, err := env.Labels.GetTable().GetAllByIndex("nameOwnerBuiltin", []interface{}{
		"Spam",
		session.Owner,
		true,
	}, []interface{}{
		"Trash",
		session.Owner,
		true,
	}, []interface{}{
		"Sent",
		session.Owner,
		true,
	}).Map(func(row r.Term) r.Term {
		return row.Field("id")
	}).Run(env.Rethink)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListUnableToFetchBuiltinLabels, err, true,
		))
		return
	}
	defer cursor.Close()
	var spamTrashSent []string
	if err := cursor.All(&spamTrashSent); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListUnableToFetchBuiltinLabels, err, true,
		))
		return
	}

	if len(spamTrashSent) != 3 {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListInvalidBuiltinLabels, "Account's builtin labels are missing", true,
		))
		return
	}

	cursor, err = env.Labels.GetTable().GetAllByIndex("owner", session.Owner).Map(func(label r.Term) r.Term {
		return env.Threads.GetTable().
			GetAllByIndex("labels", label.Field("id")).
			CoerceTo("array").
			Do(func(threads r.Term) r.Term {
			return label.Merge(map[string]interface{}{
				"total_threads_count": threads.Count(),
				"unread_threads_count": threads.Filter(func(thread r.Term) r.Term {
					return thread.Field("is_read").Not().And(
						thread.Field("labels").Map(func(label r.Term) r.Term {
							return r.Expr(spamTrashSent).Contains(label)
						}).Reduce(func(left r.Term, right r.Term) r.Term {
							return left.Or(right)
						}).Not(),
					)
				}).Count(),
			})
		})
	}).Run(env.Rethink)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListUnableToFetchAllLabels, err, true,
		))
		return
	}
	defer cursor.Close()
	var labels []*models.Label
	if err := cursor.All(&labels); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListUnableToFetchAllLabels, err, true,
		))
		return
	}

	utils.JSONResponse(w, 200, &LabelsListResponse{
		Success: true,
		Labels:  &labels,
	})
}

type LabelsCreateRequest struct {
	Name string `json:"name"`
}

// LabelsCreateResponse contains the result of the LabelsCreate request.
type LabelsCreateResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Label   *models.Label `json:"label,omitempty"`
}

// LabelsCreate does *something* - TODO
func LabelsCreate(c web.C, w http.ResponseWriter, req *http.Request) {
	// Decode the request
	var input LabelsCreateRequest
	err := utils.ParseRequest(req, &input)
	if err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.LabelsCreateInvalidInput, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the input data isn't empty
	if input.Name == "" {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.LabelsCreateInvalidInput, "Name is empty", false,
		))
		return
	}

	if _, err := env.Labels.GetLabelByNameAndOwner(session.Owner, input.Name); err == nil {
		utils.JSONResponse(w, 409, utils.NewError(
			utils.LabelsCreateAlreadyExists, "A label with such name already exists", false,
		))
		return
	}

	// Create a new label struct
	label := &models.Label{
		Resource: models.MakeResource(session.Owner, input.Name),
		Builtin:  false,
	}

	// Insert the label into the database
	if err := env.Labels.Insert(label); err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.LabelsCreateUnableToInsert, err, true,
		))
		return
	}

	utils.JSONResponse(w, 201, &LabelsCreateResponse{
		Success: true,
		Label:   label,
	})
}

// LabelsGetResponse contains the result of the LabelsGet request.
type LabelsGetResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Label   *models.Label `json:"label,omitempty"`
}

// LabelsGet does *something* - TODO
func LabelsGet(c web.C, w http.ResponseWriter, req *http.Request) {
	session := c.Env["token"].(*models.Token)

	// Fetch spam, trash and id
	cursor, err := env.Labels.GetTable().GetAllByIndex("nameOwnerBuiltin", []interface{}{
		"Spam",
		session.Owner,
		true,
	}, []interface{}{
		"Trash",
		session.Owner,
		true,
	}, []interface{}{
		"Sent",
		session.Owner,
		true,
	}).Map(func(row r.Term) r.Term {
		return row.Field("id")
	}).Run(env.Rethink)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListUnableToFetchBuiltinLabels, err, true,
		))
		return
	}
	defer cursor.Close()
	var spamTrashSent []string
	if err := cursor.All(&spamTrashSent); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsListUnableToFetchBuiltinLabels, err, true,
		))
		return
	}

	// Fetch the label
	cursor, err = env.Labels.GetTable().Get(c.URLParams["id"]).Do(func(label r.Term) r.Term {
		return env.Threads.GetTable().
			GetAllByIndex("labels", label.Field("id")).
			CoerceTo("array").
			Do(func(threads r.Term) r.Term {
			return label.Merge(map[string]interface{}{
				"total_threads_count": threads.Count(),
				"unread_threads_count": threads.Filter(func(thread r.Term) r.Term {
					return thread.Field("is_read").Not().And(
						thread.Field("labels").Map(func(label r.Term) r.Term {
							return r.Expr(spamTrashSent).Contains(label)
						}).Reduce(func(left r.Term, right r.Term) r.Term {
							return left.Or(right)
						}).Not(),
					)
				}).Count(),
			})
		})
	}).Run(env.Rethink)
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.LabelsGetUnableToGet, err, true,
		))
		return
	}
	var label *models.Label
	if err := cursor.One(&label); err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.LabelsGetUnableToGet, err, true,
		))
		return
	}

	// Check for ownership
	if label.Owner != session.Owner {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.LabelsGetNotOwned, "You're not the owner of this label", false,
		))
		return
	}

	// Write the label to the response
	utils.JSONResponse(w, 200, &LabelsGetResponse{
		Success: true,
		Label:   label,
	})
}

type LabelsUpdateRequest struct {
	Name string `json:"name"`
}

// LabelsUpdateResponse contains the result of the LabelsUpdate request.
type LabelsUpdateResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Label   *models.Label `json:"label,omitempty"`
}

// LabelsUpdate does *something* - TODO
func LabelsUpdate(c web.C, w http.ResponseWriter, req *http.Request) {
	// Decode the request
	var input LabelsUpdateRequest
	err := utils.ParseRequest(req, &input)
	if err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.LabelsUpdateInvalidInput, err, false,
		))
		return
	}

	// Get the label from the database
	label, err := env.Labels.GetLabel(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.LabelsUpdateUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if label.Owner != session.Owner {
		utils.JSONResponse(w, 403, utils.NewError(
			utils.LabelsUpdateUnableToGet, "You're not the owner of this label", false,
		))
		return
	}

	if input.Name != "" {
		label.Name = input.Name
	}

	// Perform the update
	err = env.Labels.UpdateID(c.URLParams["id"], label)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsUpdateUnableToUpdate, err, true,
		))
		return
	}

	// Write the contact to the response
	utils.JSONResponse(w, 200, &LabelsUpdateResponse{
		Success: true,
		Label:   label,
	})
}

// LabelsDeleteResponse contains the result of the LabelsDelete request.
type LabelsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LabelsDelete does *something* - TODO
func LabelsDelete(c web.C, w http.ResponseWriter, req *http.Request) {
	// Get the label from the database
	label, err := env.Labels.GetLabel(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.LabelsDeleteUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if label.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.LabelsDeleteNotOwned, "You're not the owner of this label", false,
		))
		return
	}

	// Perform the deletion
	err = env.Labels.DeleteID(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.LabelsDeleteUnableToDelete, err, true,
		))
		return
	}

	utils.JSONResponse(w, 200, &LabelsDeleteResponse{
		Success: true,
		Message: "Label successfully removed",
	})
}
