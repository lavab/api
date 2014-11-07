package routes

import (
	"net/http"

	"github.com/lavab/api/utils"
)

// LabelsListResponse contains the result of the LabelsList request.
type LabelsListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LabelsList does *something* - TODO
func LabelsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &LabelsListResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
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
