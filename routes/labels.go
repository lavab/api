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
			"error": err,
		}).Error("Unable to fetch labels")

		utils.JSONResponse(w, 500, &LabelsListResponse{
			Success: false,
			Message: "Internal error (code LA/LI/01)",
		})
		return
	}

	utils.JSONResponse(w, 200, &LabelsListResponse{
		Success: true,
		Labels:  &labels,
	})
}

// LabelsCreateResponse contains the result of the LabelsCreate request.
type LabelsCreateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LabelsCreate does *something* - TODO
func LabelsCreate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &LabelsCreateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// LabelsGetResponse contains the result of the LabelsGet request.
type LabelsGetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LabelsGet does *something* - TODO
func LabelsGet(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &LabelsGetResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// LabelsUpdateResponse contains the result of the LabelsUpdate request.
type LabelsUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LabelsUpdate does *something* - TODO
func LabelsUpdate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &LabelsUpdateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// LabelsDeleteResponse contains the result of the LabelsDelete request.
type LabelsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LabelsDelete does *something* - TODO
func LabelsDelete(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &LabelsDeleteResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
