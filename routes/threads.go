package routes

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

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
func ThreadsList(c web.C, w http.ResponseWriter, r *http.Request) {
	session := c.Env["token"].(*models.Token)

	var (
		query     = r.URL.Query()
		sortRaw   = query.Get("sort")
		offsetRaw = query.Get("offset")
		limitRaw  = query.Get("limit")
		labelsRaw = query.Get("label")
		labels    []string
		sort      []string
		offset    int
		limit     int
	)

	if offsetRaw != "" {
		o, err := strconv.Atoi(offsetRaw)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":  err,
				"offset": offset,
			}).Error("Invalid offset")

			utils.JSONResponse(w, 400, &ThreadsListResponse{
				Success: false,
				Message: "Invalid offset",
			})
			return
		}
		offset = o
	}

	if limitRaw != "" {
		l, err := strconv.Atoi(limitRaw)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"limit": limit,
			}).Error("Invalid limit")

			utils.JSONResponse(w, 400, &ThreadsListResponse{
				Success: false,
				Message: "Invalid limit",
			})
			return
		}
		limit = l
	}

	if sortRaw != "" {
		sort = strings.Split(sortRaw, ",")
	}

	if labelsRaw != "" {
		labels = strings.Split(labelsRaw, ",")
	}

	threads, err := env.Threads.List(session.Owner, sort, offset, limit, labels)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch threads")

		utils.JSONResponse(w, 500, &ThreadsListResponse{
			Success: false,
			Message: "Internal error (code TH/LI/01)",
		})
		return
	}

	if offsetRaw != "" || limitRaw != "" {
		count, err := env.Threads.CountOwnedBy(session.Owner)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to count threads")

			utils.JSONResponse(w, 500, &ThreadsListResponse{
				Success: false,
				Message: "Internal error (code TH/LI/02)",
			})
			return
		}
		w.Header().Set("X-Total-Count", strconv.Itoa(count))
	}

	utils.JSONResponse(w, 200, &ThreadsListResponse{
		Success: true,
		Threads: &threads,
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

	manifest, err := env.Emails.GetThreadManifest(thread.ID)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    thread.ID,
		}).Error("Unable to get a manifest")
	} else {
		thread.Manifest = manifest
	}

	var emails []*models.Email
	if ok := r.URL.Query().Get("list_emails"); ok == "true" || ok == "1" {
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

type ThreadsUpdateRequest struct {
	Labels   []string `json:"labels"`
	IsRead   *bool    `json:"is_read"`
	LastRead *string  `json:"last_read"`
}

type ThreadsUpdateResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Thread  *models.Thread `json:"thread,omitempty"`
}

// ThreadsUpdate does *something* with a thread.
func ThreadsUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	var input ThreadsUpdateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &ThreadsUpdateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the thread from the database
	thread, err := env.Threads.GetThread(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &ThreadsUpdateResponse{
			Success: false,
			Message: "Thread not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if thread.Owner != session.Owner {
		utils.JSONResponse(w, 404, &ContactsUpdateResponse{
			Success: false,
			Message: "Contact not found",
		})
		return
	}

	if input.Labels != nil && !reflect.DeepEqual(thread.Labels, input.Labels) {
		thread.Labels = input.Labels
	}

	if input.LastRead != nil && *input.LastRead != thread.LastRead {
		thread.LastRead = *input.LastRead
	}

	if input.IsRead != nil && *input.IsRead != thread.IsRead {
		thread.IsRead = *input.IsRead
	}

	// Disabled for now, as we're using DateModified for sorting by the date of the last email
	// thread.DateModified = time.Now()

	err = env.Threads.UpdateID(c.URLParams["id"], thread)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to update a thread")

		utils.JSONResponse(w, 500, &ThreadsUpdateResponse{
			Success: false,
			Message: "Internal error (code TH/UP/01)",
		})
		return
	}

	// Write the thread to the response
	utils.JSONResponse(w, 200, &ThreadsUpdateResponse{
		Success: true,
		Thread:  thread,
	})
}

// ThreadsDeleteResponse contains the result of the ThreadsDelete request.
type ThreadsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ThreadsDelete removes a thread from the database
func ThreadsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the thread from the database
	thread, err := env.Threads.GetThread(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &ThreadsDeleteResponse{
			Success: false,
			Message: "Thread not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if thread.Owner != session.Owner {
		utils.JSONResponse(w, 404, &ThreadsDeleteResponse{
			Success: false,
			Message: "Thread not found",
		})
		return
	}

	// Perform the deletion
	err = env.Threads.DeleteID(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete a thread")

		utils.JSONResponse(w, 500, &ThreadsDeleteResponse{
			Success: false,
			Message: "Internal error (code TH/DE/01)",
		})
		return
	}

	// Remove dependent emails
	err = env.Emails.DeleteByThread(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete emails by thread")

		utils.JSONResponse(w, 500, &ThreadsDeleteResponse{
			Success: false,
			Message: "Internal error (code TH/DE/02)",
		})
		return
	}

	// Write the thread to the response
	utils.JSONResponse(w, 200, &ThreadsDeleteResponse{
		Success: true,
		Message: "Thread successfully removed",
	})
}
