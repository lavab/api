package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// ThreadsListResponse contains the result of the ThreadsList request.
type ThreadsListResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Threads *[]*models.Thread `json:"threads,omitempty"`
}

// ThreadsList shows all threads
func ThreadsList(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ThreadsListResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}

// ThreadsGetResponse contains the result of the ThreadsGet request.
type ThreadsGetResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Thread  *models.Thread   `json:"thread,omitempty"`
	Emails  *[]*models.Email `json:"emails,omitempty"`
}

// ThreadsGet returns information about a single thread.
func ThreadsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	thread, err := env.Threads.GetThread(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &ThreadsGetResponse{
			Success: false,
			Message: "Thread not found",
		})
		return
	}

	session := c.Env["token"].(*models.Token)

	if thread.Owner != session.Owner {
		utils.JSONResponse(w, 404, &ThreadsGetResponse{
			Success: false,
			Message: "Thread not found",
		})
		return
	}

	var emails []*models.Email
	if ok := r.URL.Query().Get("list_emails"); ok == "true" {
		emails, err = env.Emails.GetByThread(thread.ID)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    thread.ID,
			}).Error("Unable to fetch emails linked to a thread")

			utils.JSONResponse(w, 500, &ThreadsGetResponse{
				Success: false,
				Message: "Unable to retrieve emails",
			})
			return
		}
	}

	utils.JSONResponse(w, 200, &ThreadsGetResponse{
		Success: true,
		Thread:  thread,
		Emails:  &emails,
	})
}

// ThreadsUpdateResponse contains the result of the ThreadsUpdate request.
type ThreadsUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ThreadsUpdate does *something* with a thread.
func ThreadsUpdate(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &ThreadsUpdateResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
