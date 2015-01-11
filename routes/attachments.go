package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
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
			"error": err,
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
