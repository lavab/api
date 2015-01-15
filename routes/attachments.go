package routes

import (
	"net/http"

	"github.com/lavab/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

type AttachmentsListResponse struct {
	Success     bool                  `json:"success"`
	Message     string                `json:"message,omitempty"`
	Attachments *[]*models.Attachment `json:"attachments,omitempty"`
}

func AttachmentsList(c web.C, w http.ResponseWriter, r *http.Request) {
	session := c.Env["token"].(*models.Token)

	attachments, err := env.Attachments.GetOwnedBy(session.Owner)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch attachments")

		utils.JSONResponse(w, 500, &AttachmentsListResponse{
			Success: false,
			Message: "Internal error (code AT/LI/01)",
		})
		return
	}

	utils.JSONResponse(w, 200, &AttachmentsListResponse{
		Success:     true,
		Attachments: &attachments,
	})
}

type AttachmentsCreateRequest struct {
	Data            string   `json:"data" schema:"data"`
	Name            string   `json:"name" schema:"name"`
	Encoding        string   `json:"encoding" schema:"encoding"`
	VersionMajor    int      `json:"version_major" schema:"version_major"`
	VersionMinor    int      `json:"version_minor" schema:"version_minor"`
	PGPFingerprints []string `json:"pgp_fingerprints" schema:"pgp_fingerprints"`
}

type AttachmentsCreateResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Attachment *models.Attachment `json:"attachment,omitempty"`
}

// AttachmentsCreate creates a new attachment
func AttachmentsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input AttachmentsCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &AttachmentsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the input data isn't empty
	if input.Data == "" || input.Name == "" || input.Encoding == "" ||
		input.PGPFingerprints == nil || len(input.PGPFingerprints) == 0 {
		utils.JSONResponse(w, 400, &AttachmentsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Create a new attachment struct
	attachment := &models.Attachment{
		Encrypted: models.Encrypted{
			Encoding:        input.Encoding,
			Data:            input.Data,
			Schema:          "attachment",
			VersionMajor:    input.VersionMajor,
			VersionMinor:    input.VersionMinor,
			PGPFingerprints: input.PGPFingerprints,
		},
		Resource: models.MakeResource(session.Owner, input.Name),
	}

	// Insert the attachment into the database
	if err := env.Attachments.Insert(attachment); err != nil {
		utils.JSONResponse(w, 500, &AttachmentsCreateResponse{
			Success: false,
			Message: "internal server error - AT/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not insert a attachment into the database")
		return
	}

	utils.JSONResponse(w, 201, &AttachmentsCreateResponse{
		Success:    true,
		Message:    "A new attachment was successfully created",
		Attachment: attachment,
	})
}

// AttachmentsGetResponse contains the result of the AttachmentsGet request.
type AttachmentsGetResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message,omitempty"`
	Attachment *models.Attachment `json:"attachment,omitempty"`
}

// AttachmentsGet gets the requested attachment from the database
func AttachmentsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the attachment from the database
	attachment, err := env.Attachments.GetAttachment(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &AttachmentsGetResponse{
			Success: false,
			Message: "Attachment not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if attachment.Owner != session.Owner {
		utils.JSONResponse(w, 404, &AttachmentsGetResponse{
			Success: false,
			Message: "Attachment not found",
		})
		return
	}

	// Write the attachment to the response
	utils.JSONResponse(w, 200, &AttachmentsGetResponse{
		Success:    true,
		Attachment: attachment,
	})
}

// AttachmentsUpdateRequest is the payload passed to PUT /contacts/:id
type AttachmentsUpdateRequest struct {
	Data            string   `json:"data" schema:"data"`
	Name            string   `json:"name" schema:"name"`
	Encoding        string   `json:"encoding" schema:"encoding"`
	VersionMajor    *int     `json:"version_major" schema:"version_major"`
	VersionMinor    *int     `json:"version_minor" schema:"version_minor"`
	PGPFingerprints []string `json:"pgp_fingerprints" schema:"pgp_fingerprints"`
}

// AttachmentsUpdateResponse contains the result of the AttachmentsUpdate request.
type AttachmentsUpdateResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message,omitempty"`
	Attachment *models.Attachment `json:"attachment,omitempty"`
}

// AttachmentsUpdate updates an existing attachment in the database
func AttachmentsUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input AttachmentsUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &AttachmentsUpdateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the attachment from the database
	attachment, err := env.Attachments.GetAttachment(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &AttachmentsUpdateResponse{
			Success: false,
			Message: "Attachment not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if attachment.Owner != session.Owner {
		utils.JSONResponse(w, 404, &AttachmentsUpdateResponse{
			Success: false,
			Message: "Attachment not found",
		})
		return
	}

	if input.Data != "" {
		attachment.Data = input.Data
	}

	if input.Name != "" {
		attachment.Name = input.Name
	}

	if input.Encoding != "" {
		attachment.Encoding = input.Encoding
	}

	if input.VersionMajor != nil {
		attachment.VersionMajor = *input.VersionMajor
	}

	if input.VersionMinor != nil {
		attachment.VersionMinor = *input.VersionMinor
	}

	if input.PGPFingerprints != nil {
		attachment.PGPFingerprints = input.PGPFingerprints
	}

	// Perform the update
	err = env.Attachments.UpdateID(c.URLParams["id"], input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to update a attachment")

		utils.JSONResponse(w, 500, &AttachmentsUpdateResponse{
			Success: false,
			Message: "Internal error (code AT/UP/01)",
		})
		return
	}

	// Write the attachment to the response
	utils.JSONResponse(w, 200, &AttachmentsUpdateResponse{
		Success:    true,
		Attachment: attachment,
	})
}

// AttachmentsDeleteResponse contains the result of the Delete request.
type AttachmentsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AttachmentsDelete removes a attachment from the database
func AttachmentsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the attachment from the database
	attachment, err := env.Attachments.GetAttachment(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &AttachmentsDeleteResponse{
			Success: false,
			Message: "Attachment not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if attachment.Owner != session.Owner {
		utils.JSONResponse(w, 404, &AttachmentsDeleteResponse{
			Success: false,
			Message: "Attachment not found",
		})
		return
	}

	// Perform the deletion
	err = env.Attachments.DeleteID(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete a attachment")

		utils.JSONResponse(w, 500, &AttachmentsDeleteResponse{
			Success: false,
			Message: "Internal error (code AT/DE/01)",
		})
		return
	}

	// Write the attachment to the response
	utils.JSONResponse(w, 200, &AttachmentsDeleteResponse{
		Success: true,
		Message: "Attachment successfully removed",
	})
}
