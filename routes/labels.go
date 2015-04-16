package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
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
func LabelsList(c web.C, w http.ResponseWriter, r *http.Request) {
	session := c.Env["token"].(*models.Token)

	labels, err := env.Labels.GetOwnedBy(session.Owner)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch labels")

		utils.JSONResponse(w, 500, &LabelsListResponse{
			Success: false,
			Message: "Internal error (code LA/LI/01)",
		})
		return
	}

	for _, label := range labels {
		if label.Builtin && label.Name != "Inbox" {
			continue
		}

		totalCount, err := env.Threads.CountByLabel(label.ID)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"label": label.ID,
			}).Error("Unable to fetch total threads count")

			utils.JSONResponse(w, 500, &LabelsListResponse{
				Success: false,
				Message: "Internal error (code LA/LI/02)",
			})
			return
		}

		unreadCount, err := env.Threads.CountByLabelUnread(label.ID)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"label": label.ID,
			}).Error("Unable to fetch unread threads count")

			utils.JSONResponse(w, 500, &LabelsListResponse{
				Success: false,
				Message: "Internal error (code LA/LI/03)",
			})
			return
		}

		label.TotalThreadsCount = totalCount
		label.UnreadThreadsCount = unreadCount
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
func LabelsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input LabelsCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &LabelsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the input data isn't empty
	if input.Name == "" {
		utils.JSONResponse(w, 400, &LabelsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if _, err := env.Labels.GetLabelByNameAndOwner(session.Owner, input.Name); err == nil {
		utils.JSONResponse(w, 409, &LabelsCreateResponse{
			Success: false,
			Message: "Label with such name already exists",
		})
		return
	}

	// Create a new label struct
	label := &models.Label{
		Resource: models.MakeResource(session.Owner, input.Name),
		Builtin:  false,
	}

	// Insert the label into the database
	if err := env.Labels.Insert(label); err != nil {
		utils.JSONResponse(w, 500, &LabelsCreateResponse{
			Success: false,
			Message: "internal server error - LA/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not insert a label into the database")
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
func LabelsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the label from the database
	label, err := env.Labels.GetLabel(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &LabelsGetResponse{
			Success: false,
			Message: "Label not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if label.Owner != session.Owner {
		utils.JSONResponse(w, 404, &LabelsGetResponse{
			Success: false,
			Message: "Label not found",
		})
		return
	}

	totalCount, err := env.Threads.CountByLabel(label.ID)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"label": label.ID,
		}).Error("Unable to fetch total threads count")

		utils.JSONResponse(w, 500, &LabelsListResponse{
			Success: false,
			Message: "Internal error (code LA/GE/01)",
		})
		return
	}

	unreadCount, err := env.Threads.CountByLabelUnread(label.ID)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"label": label.ID,
		}).Error("Unable to fetch unread threads count")

		utils.JSONResponse(w, 500, &LabelsListResponse{
			Success: false,
			Message: "Internal error (code LA/GE/01)",
		})
		return
	}

	label.TotalThreadsCount = totalCount
	label.UnreadThreadsCount = unreadCount

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
func LabelsUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input LabelsUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &LabelsUpdateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the label from the database
	label, err := env.Labels.GetLabel(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &LabelsUpdateResponse{
			Success: false,
			Message: "Label not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if label.Owner != session.Owner {
		utils.JSONResponse(w, 404, &LabelsUpdateResponse{
			Success: false,
			Message: "Label not found",
		})
		return
	}

	if input.Name != "" {
		label.Name = input.Name
	}

	// Perform the update
	err = env.Labels.UpdateID(c.URLParams["id"], label)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to update a contact")

		utils.JSONResponse(w, 500, &LabelsUpdateResponse{
			Success: false,
			Message: "Internal error (code LA/UP/01)",
		})
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
func LabelsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the label from the database
	label, err := env.Labels.GetLabel(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &LabelsDeleteResponse{
			Success: false,
			Message: "Label not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if label.Owner != session.Owner {
		utils.JSONResponse(w, 404, &LabelsDeleteResponse{
			Success: false,
			Message: "Label not found",
		})
		return
	}

	// Perform the deletion
	err = env.Labels.DeleteID(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete a label")

		utils.JSONResponse(w, 500, &LabelsDeleteResponse{
			Success: false,
			Message: "Internal error (code LA/DE/01)",
		})
		return
	}

	utils.JSONResponse(w, 200, &LabelsDeleteResponse{
		Success: true,
		Message: "Label successfully removed",
	})
}
