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
			utils.JSONResponse(w, 400, utils.NewError(
				utils.ThreadsListInvalidOffset, err, false,
			))
			return
		}
		offset = o
	}

	if limitRaw != "" {
		l, err := strconv.Atoi(limitRaw)
		if err != nil {
			utils.JSONResponse(w, 400, utils.NewError(
				utils.ThreadsListInvalidLimit, err, false,
			))
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
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ThreadsListUnableToGet, err, true,
		))
		return
	}

	if offsetRaw != "" || limitRaw != "" {
		count, err := env.Threads.CountOwnedBy(session.Owner)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.ThreadsListUnableToCount, err, true,
			))
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
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Thread  *models.Thread  `json:"thread,omitempty"`
	Emails  []*models.Email `json:"emails,omitempty"`
}

// ThreadsGet returns information about a single thread.
func ThreadsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	thread, err := env.Threads.GetThread(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ThreadsGetUnableToGet, err, false,
		))
		return
	}

	session := c.Env["token"].(*models.Token)

	if thread.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ThreadsGetNotOwned, "You're not the owner of this thread", false,
		))
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
			utils.JSONResponse(w, 500, utils.NewError(
				utils.ThreadsGetUnableToFetchEmails, err, true,
			))
			return
		}
	}

	utils.JSONResponse(w, 200, &ThreadsGetResponse{
		Success: true,
		Thread:  thread,
		Emails:  emails,
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
		utils.JSONResponse(w, 400, utils.NewError(
			utils.ThreadsUpdateInvalidInput, err, false,
		))
		return
	}

	// Get the thread from the database
	thread, err := env.Threads.GetThread(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ThreadsUpdateUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if thread.Owner != session.Owner {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.ThreadsUpdateNotOwned, "You're not the owner of this thread", false,
		))
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
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ThreadsUpdateUnableToUpdate, err, true,
		))
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
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ThreadsUpdateUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if thread.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.ThreadsUpdateNotOwned, "You're not the owner of this thread", false,
		))
		return
	}

	// Perform the deletion
	err = env.Threads.DeleteID(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ThreadsDeleteUnableToDeleteThread, err, true,
		))
		return
	}

	// Remove dependent emails
	err = env.Emails.DeleteByThread(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.ThreadsDeleteUnableToDeleteEmails, err, true,
		))
		return
	}

	// Write the thread to the response
	utils.JSONResponse(w, 200, &ThreadsDeleteResponse{
		Success: true,
		Message: "Thread successfully removed",
	})
}
