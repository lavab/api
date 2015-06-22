package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/zenazn/goji/web"
	_ "golang.org/x/crypto/ripemd160"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

var prefixesRegex = regexp.MustCompile(`([\[\(] *)?(RE?S?|FYI|RIF|I|FS|VB|RV|ENC|ODP|PD|YNT|ILT|SV|VS|VL|AW|WG|ΑΠ|ΣΧΕΤ|ΠΡΘ|תגובה|הועבר|主题|转发|FWD?) *([-:;)\]][ :;\])-]*|$)|\]+ *$`)

// EmailsListResponse contains the result of the EmailsList request.
type EmailsListResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Emails  *[]*models.Email `json:"emails,omitempty"`
}

// EmailsList sends a list of the emails in the inbox.
func EmailsList(c web.C, w http.ResponseWriter, r *http.Request) {
	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Parse the query
	var (
		query     = r.URL.Query()
		sortRaw   = query.Get("sort")
		offsetRaw = query.Get("offset")
		limitRaw  = query.Get("limit")
		thread    = query.Get("thread")
		sort      []string
		offset    int
		limit     int
	)

	if offsetRaw != "" {
		o, err := strconv.Atoi(offsetRaw)
		if err != nil {
			utils.JSONResponse(w, 400, utils.NewError(
				utils.EmailsListInvalidOffset, err, false,
			))
			return
		}
		offset = o
	}

	if limitRaw != "" {
		l, err := strconv.Atoi(limitRaw)
		if err != nil {
			utils.JSONResponse(w, 400, utils.NewError(
				utils.EmailsListInvalidLimit, err, false,
			))
			return
		}
		limit = l
	}

	if sortRaw != "" {
		sort = strings.Split(sortRaw, ",")
	}

	// Get contacts from the database
	emails, err := env.Emails.List(session.Owner, sort, offset, limit, thread)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.EmailsListUnableToGet, err, true,
		))
		return
	}

	if offsetRaw != "" || limitRaw != "" {
		count, err := env.Emails.CountOwnedBy(session.Owner)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.EmailsListUnableToCount, err, true,
			))
			return
		}
		w.Header().Set("X-Total-Count", strconv.Itoa(count))
	}

	utils.JSONResponse(w, 200, &EmailsListResponse{
		Success: true,
		Emails:  &emails,
	})

	// GET parameters:
	//   sort - split by commas, prefixes: - is desc, + is asc
	//   offset, limit - for pagination
	// Pagination ADDS X-Total-Count to the response!
}

type EmailsCreateRequest struct {
	// Internal properties
	Kind   string `json:"kind"`
	Thread string `json:"thread"`

	// Metadata that has to be leaked
	From string   `json:"from"`
	To   []string `json:"to"`
	CC   []string `json:"cc"`
	BCC  []string `json:"bcc"`

	// Encrypted parts
	Manifest string   `json:"manifest"`
	Body     string   `json:"body"`
	Files    []string `json:"files"`

	// Temporary partials if you're sending unencrypted
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	InReplyTo   string `json:"in_reply_to"`

	SubjectHash string `json:"subject_hash"`
}

// EmailsCreateResponse contains the result of the EmailsCreate request.
type EmailsCreateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Created []string `json:"created,omitempty"`
}

// EmailsCreate sends a new email
func EmailsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input EmailsCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.EmailsCreateInvalidInput, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the kind is valid
	if input.Kind != "raw" && input.Kind != "manifest" && input.Kind != "pgpmime" {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.EmailsCreateInvalidInput, "Invalid email encryption kind", false,
		))
		return
	}

	// Ensure that there's at least one recipient and that there's body
	if input.To == nil || len(input.To) == 0 || input.Body == "" {
		utils.JSONResponse(w, 400, utils.NewError(
			utils.EmailsCreateInvalidInput, "Invalid to field or empty body", false,
		))
		return
	}

	if input.Files != nil && len(input.Files) > 0 {
		// Check rights to files
		files, err := env.Files.GetFiles(input.Files...)
		if err != nil {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.EmailsCreateUnableToFetchFiles, err, false,
			))
			return
		}
		for _, file := range files {
			if file.Owner != session.Owner {
				utils.JSONResponse(w, 403, utils.NewError(
					utils.EmailsCreateFileNotOwned, "You're not the owner of file "+file.ID, false,
				))
				return
			}
		}
	}

	// Create an email resource
	resource := models.MakeResource(session.Owner, input.Subject)

	// Generate metadata for manifests
	if input.Kind == "manifest" {
		resource.Name = "Encrypted message (" + resource.ID + ")"
	}

	// Fetch the user object from the database
	account, err := env.Accounts.GetTokenOwner(c.Env["token"].(*models.Token))
	if err != nil {
		// The session refers to a non-existing user
		utils.JSONResponse(w, 500, utils.NewError(
			utils.EmailsCreateUnableToFetchAccount, err, true,
		))
		return
	}

	// Get the "Sent" label's ID
	var label *models.Label
	err = env.Labels.WhereAndFetchOne(map[string]interface{}{
		"name":    "Sent",
		"builtin": true,
		"owner":   account.ID,
	}, &label)
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.EmailsCreateUnableToFetchLabel, err, true,
		))
		return
	}

	if input.From != "" {
		// Parse the from field
		from, err := mail.ParseAddress(input.From)
		if err != nil {
			utils.JSONResponse(w, 400, utils.NewError(
				utils.EmailsCreateInvalidFromAddress, err, false,
			))
			return
		}

		// We have a specified address
		if from.Address != "" {
			parts := strings.SplitN(from.Address, "@", 2)

			if parts[1] != env.Config.EmailDomain {
				utils.JSONResponse(w, 403, utils.NewError(
					utils.EmailsCreateInvalidFromAddress, "Invalid from domain", false,
				))
				return
			}

			address, err := env.Addresses.GetAddress(parts[0])
			if err != nil {
				utils.JSONResponse(w, 403, utils.NewError(
					utils.EmailsCreateInvalidFromAddress, err, false,
				))
				return
			}

			if address.Owner != account.ID {
				utils.JSONResponse(w, 403, utils.NewError(
					utils.EmailsCreateInvalidFromAddress, "You're not the owner of that address.", false,
				))
				return
			}
		}
	} else {
		displayName := ""

		if x, ok := account.Settings.(map[string]interface{}); ok {
			if y, ok := x["displayName"]; ok {
				if z, ok := y.(string); ok {
					displayName = z
				}
			}
		}

		addr := &mail.Address{
			Name:    displayName,
			Address: account.StyledName + "@" + env.Config.EmailDomain,
		}

		input.From = addr.String()
	}

	// Check if Thread is set
	if input.Thread != "" {
		// todo: make it an actual exists check to reduce lan bandwidth
		thread, err := env.Threads.GetThread(input.Thread)
		if err != nil {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.EmailsCreateUnableToFetchThread, err, false,
			))
			return
		}

		if thread.Owner != account.Owner {
			utils.JSONResponse(w, 403, utils.NewError(
				utils.EmailsCreateThreadNotOwned, "You're not the owner of this thread", false,
			))
			return
		}

		// update thread.secure depending on email's kind
		if (input.Kind == "raw" && thread.Secure == "all") ||
			(input.Kind == "manifest" && thread.Secure == "none") ||
			(input.Kind == "pgpmime" && thread.Secure == "none") {
			if err := env.Threads.UpdateID(thread.ID, map[string]interface{}{
				"secure": "some",
			}); err != nil {
				utils.JSONResponse(w, 500, utils.NewError(
					utils.EmailsCreateUnableToUpdateThread, err, true,
				))
				return
			}
		}
	} else {
		secure := "all"
		if input.Kind == "raw" {
			secure = "none"
		}

		if input.SubjectHash == "" {
			// Generate the subject hash
			shr := sha256.Sum256([]byte(prefixesRegex.ReplaceAllString(input.Subject, "")))
			input.SubjectHash = hex.EncodeToString(shr[:])
		}

		thread := &models.Thread{
			Resource:    models.MakeResource(account.ID, "Encrypted thread"),
			Emails:      []string{resource.ID},
			Labels:      []string{label.ID},
			Members:     append(append(input.To, input.CC...), input.BCC...),
			IsRead:      true,
			SubjectHash: input.SubjectHash,
			Secure:      secure,
		}

		err := env.Threads.Insert(thread)
		if err != nil {
			utils.JSONResponse(w, 500, utils.NewError(
				utils.EmailsCreateUnableToInsertThread, err, true,
			))
			return
		}

		input.Thread = thread.ID
	}

	// Calculate the message ID
	idHash := sha256.Sum256([]byte(resource.ID))
	messageID := hex.EncodeToString(idHash[:]) + "@" + env.Config.EmailDomain

	// Determine if email is secure
	secure := true
	if input.Kind == "raw" {
		secure = false
	}

	// Create a new email struct
	email := &models.Email{
		Resource:  resource,
		MessageID: messageID,

		Kind:   input.Kind,
		Thread: input.Thread,

		From: input.From,
		To:   input.To,
		CC:   input.CC,
		BCC:  input.BCC,

		Manifest: input.Manifest,
		Body:     input.Body,
		Files:    input.Files,

		ContentType: input.ContentType,
		InReplyTo:   input.InReplyTo,

		Status: "queued",
		Secure: secure,
	}

	// Insert the email into the database
	if err := env.Emails.Insert(email); err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.EmailsCreateUnableToInsertEmail, err, true,
		))
		return
	}

	// Add a send request to the queue
	err = env.Producer.Publish("send_email", []byte(`"`+email.ID+`"`))
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.EmailsCreateUnableToQueue, err, true,
		))
		return
	}

	utils.JSONResponse(w, 201, &EmailsCreateResponse{
		Success: true,
		Created: []string{email.ID},
	})
}

// EmailsGetResponse contains the result of the EmailsGet request.
type EmailsGetResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Email   *models.Email `json:"email,omitempty"`
}

// EmailsGet responds with a single email message
func EmailsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the email from the database
	email, err := env.Emails.GetEmail(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.EmailsGetUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if email.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.EmailsGetNotOwned, "You're not the owner of this email", false,
		))
		return
	}

	// Write the email to the response
	utils.JSONResponse(w, 200, &EmailsGetResponse{
		Success: true,
		Email:   email,
	})
}

// EmailsDeleteResponse contains the result of the EmailsDelete request.
type EmailsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EmailsDelete remvoes an email from the system
func EmailsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the email from the database
	email, err := env.Emails.GetEmail(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.EmailsDeleteUnableToGet, err, false,
		))
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if email.Owner != session.Owner {
		utils.JSONResponse(w, 404, utils.NewError(
			utils.EmailsDeleteNotOwned, "You're not the owner of this email", false,
		))
		return
	}

	// Perform the deletion
	err = env.Emails.DeleteID(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 500, utils.NewError(
			utils.EmailsDeleteUnableToDelete, err, false,
		))
		return
	}

	// Write the email to the response
	utils.JSONResponse(w, 200, &EmailsDeleteResponse{
		Success: true,
		Message: "Email successfully removed",
	})
}
