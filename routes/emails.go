package routes

import (
	//"bytes"
	//"io"
	//"crypto/sha256"
	//"encoding/hex"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	//"golang.org/x/crypto/openpgp"
	//"golang.org/x/crypto/openpgp/armor"
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
			env.Log.WithFields(logrus.Fields{
				"error":  err,
				"offset": offset,
			}).Error("Invalid offset")

			utils.JSONResponse(w, 400, &EmailsListResponse{
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

			utils.JSONResponse(w, 400, &EmailsListResponse{
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

	// Get contacts from the database
	emails, err := env.Emails.List(session.Owner, sort, offset, limit, thread)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch emails")

		utils.JSONResponse(w, 500, &EmailsListResponse{
			Success: false,
			Message: "Internal error (code EM/LI/01)",
		})
		return
	}

	if offsetRaw != "" || limitRaw != "" {
		count, err := env.Emails.CountOwnedBy(session.Owner)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to count emails")

			utils.JSONResponse(w, 500, &EmailsListResponse{
				Success: false,
				Message: "Internal error (code EM/LI/02)",
			})
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
	To  []string `json:"to"`
	CC  []string `json:"cc"`
	BCC []string `json:"bcc"`

	// Encrypted parts
	PGPFingerprints []string `json:"pgp_fingerprints"`
	Manifest        string   `json:"manifest"`
	Body            string   `json:"body"`
	Files           []string `json:"files"`

	// Temporary partials if you're sending unencrypted
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	ReplyTo     string `json:"reply_to"`

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
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &EmailsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the kind is valid
	if input.Kind != "raw" && input.Kind != "manifest" && input.Kind != "pgpmime" {
		utils.JSONResponse(w, 400, &EmailsCreateResponse{
			Success: false,
			Message: "Invalid email encryption kind",
		})
		return
	}

	// Ensure that there's at least one recipient and that there's body
	if len(input.To) == 0 || input.Body == "" {
		utils.JSONResponse(w, 400, &EmailsCreateResponse{
			Success: false,
			Message: "Invalid email",
		})
		return
	}

	if input.Files != nil && len(input.Files) > 0 {
		// Check rights to files
		files, err := env.Files.GetFiles(input.Files...)
		if err != nil {
			utils.JSONResponse(w, 500, &EmailsCreateResponse{
				Success: false,
				Message: "Unable to fetch emails",
			})
			return
		}
		for _, file := range files {
			if file.Owner != session.Owner {
				utils.JSONResponse(w, 403, &EmailsCreateResponse{
					Success: false,
					Message: "You are not the owner of file " + file.ID,
				})
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
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err.Error(),
		}).Warn("Valid session referred to a removed account")

		utils.JSONResponse(w, 410, &EmailsCreateResponse{
			Success: false,
			Message: "Account disabled",
		})
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
		env.Log.WithFields(logrus.Fields{
			"id":    account.ID,
			"error": err.Error(),
		}).Warn("Account has no sent label")

		utils.JSONResponse(w, 410, &EmailsCreateResponse{
			Success: false,
			Message: "Misconfigured account",
		})
		return
	}

	// Check if Thread is set
	if input.Thread != "" {
		// todo: make it an actual exists check to reduce lan bandwidth
		thread, err := env.Threads.GetThread(input.Thread)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    input.Thread,
				"error": err.Error(),
			}).Warn("Cannot retrieve a thread")

			utils.JSONResponse(w, 400, &EmailsCreateResponse{
				Success: false,
				Message: "Invalid thread",
			})
			return
		}

		// update thread.secure depending on email's kind
		if (input.Kind == "raw" && thread.Secure == "all") ||
			(input.Kind == "manifest" && thread.Secure == "none") ||
			(input.Kind == "pgpmime" && thread.Secure == "none") {
			if err := env.Threads.UpdateID(thread.ID, map[string]interface{}{
				"secure": "some",
			}); err != nil {
				env.Log.WithFields(logrus.Fields{
					"id":    input.Thread,
					"error": err.Error(),
				}).Warn("Cannot update a thread")

				utils.JSONResponse(w, 400, &EmailsCreateResponse{
					Success: false,
					Message: "Unable to update the thread",
				})
				return
			}
		}
	} else {
		secure := "all"
		if input.Kind == "raw" {
			secure = "none"
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
			utils.JSONResponse(w, 500, &EmailsCreateResponse{
				Success: false,
				Message: "Unable to create a new thread",
			})

			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to create a new thread")
			return
		}

		input.Thread = thread.ID
	}

	// Create a new email struct
	email := &models.Email{
		Resource:  resource,
		MessageID: resource.ID + "@lavaboom.com",

		Kind:   input.Kind,
		Thread: input.Thread,

		From: account.StyledName + "@" + env.Config.EmailDomain,
		To:   input.To,
		CC:   input.CC,
		BCC:  input.BCC,

		PGPFingerprints: input.PGPFingerprints,
		Manifest:        input.Manifest,
		Body:            input.Body,
		Files:           input.Files,

		ContentType: input.ContentType,
		ReplyTo:     input.ReplyTo,

		Status: "queued",
	}

	// Insert the email into the database
	if err := env.Emails.Insert(email); err != nil {
		utils.JSONResponse(w, 500, &EmailsCreateResponse{
			Success: false,
			Message: "internal server error - EM/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not insert an email into the database")
		return
	}

	// I'm going to whine at this part, as we are doubling the email sending code

	// Check if To contains lavaboom emails
	/*for _, address := range email.To {
		parts := strings.SplitN(address, "@", 2)
		if parts[1] == env.Config.EmailDomain {
			go sendEmail(parts[0], email)
		}
	}

	// Check if CC contains lavaboom emails
	for _, address := range email.CC {
		parts := strings.SplitN(address, "@", 2)
		if parts[1] == env.Config.EmailDomain {
			go sendEmail(parts[0], email)
		}
	}

	// Check if BCC contains lavaboom emails
	for _, address := range email.BCC {
		parts := strings.SplitN(address, "@", 2)
		if parts[1] == env.Config.EmailDomain {
			go sendEmail(parts[0], email)
		}
	}*/

	// Add a send request to the queue
	err = env.Producer.Publish("send_email", []byte(`"`+email.ID+`"`))
	if err != nil {
		utils.JSONResponse(w, 500, &EmailsCreateResponse{
			Success: false,
			Message: "internal server error - EM/CR/03",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not publish an email send request")
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
		utils.JSONResponse(w, 404, &EmailsGetResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if email.Owner != session.Owner {
		utils.JSONResponse(w, 404, &EmailsGetResponse{
			Success: false,
			Message: "Email not found",
		})
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
		utils.JSONResponse(w, 404, &EmailsDeleteResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if email.Owner != session.Owner {
		utils.JSONResponse(w, 404, &EmailsDeleteResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Perform the deletion
	err = env.Emails.DeleteID(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete a email")

		utils.JSONResponse(w, 500, &EmailsDeleteResponse{
			Success: false,
			Message: "Internal error (code EM/DE/01)",
		})
		return
	}

	// Write the email to the response
	utils.JSONResponse(w, 200, &EmailsDeleteResponse{
		Success: true,
		Message: "Email successfully removed",
	})
}

/*func sendEmail(account string, email *models.Email) {
	// find recipient's account
	recipient, err := env.Accounts.FindAccountByName(account)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"name":  account,
		}).Warn("Unable to fetch recipent's account")
		return
	}

	newEmail := *email

	// check if the email is unencrypted
	if newEmail.Body.PGPFingerprints == nil || len(newEmail.Body.PGPFingerprints) == 0 {
		// check if the acc has a pkey set
		if recipient.PublicKey == "" {
			env.Log.WithFields(logrus.Fields{
				"name": account,
			}).Warn("Recipient has no public key set")
			return
		}

		// fetch the pkey
		key, err := env.Keys.FindByFingerprint(recipient.PublicKey)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": recipient.PublicKey,
				"name":        account,
			}).Warn("Recipient's public key does not exist")
			return
		}

		// parse the armored key
		entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(key.Key))
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": recipient.PublicKey,
			}).Warn("Cannot parse an armored key")
			return
		}

		// first key should be the pkey
		publicKey := entityList[0]

		// prepare a buffer for ciphertext and initialize openpgp
		output := &bytes.Buffer{}
		input, err := openpgp.Encrypt(output, []*openpgp.Entity{publicKey}, nil, nil, nil)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": recipient.PublicKey,
			}).Warn("Cannot set up an OpenPGP encrypter")
			return
		}

		// write email's contents into input
		_, err = input.Write([]byte(newEmail.Body.Data))
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": recipient.PublicKey,
			}).Warn("Cannot write into the OpenPGP's input")
			return
		}

		// close the input
		if err := input.Close(); err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": recipient.PublicKey,
			}).Warn("Cannot close OpenPGP's input")
			return
		}

		// encode output into armor
		armoredOutput := &bytes.Buffer{}
		armoredInput, err := armor.Encode(armoredOutput, "PGP MESSAGE", map[string]string{
			"Version": "Lavaboom " + env.Config.APIVersion,
		})
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Cannot initialize a new armor encoding")
			return
		}

		_, err = io.Copy(armoredInput, output)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Unable to copy encrypted ciphertext into the armor processor")
			return
		}

		if err := armoredInput.Close(); err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"fingerprint": recipient.PublicKey,
			}).Warn("Cannot close armoring's input")
			return
		}

		newEmail.Body.PGPFingerprints = []string{recipient.PublicKey}
		newEmail.Body.Data = armoredOutput.String()
	}

	// Get the "Inbox" label's ID
	var inbox *models.Label
	err = env.Labels.WhereAndFetchOne(map[string]interface{}{
		"name":    "Inbox",
		"builtin": true,
		"owner":   recipient.ID,
	}, &inbox)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"id":    recipient.ID,
			"error": err.Error(),
		}).Warn("Account has no inbox label")
		return
	}

	// strip prefixes from the subject
	rawSubject := prefixesRegex.ReplaceAllString(newEmail.Name, "")

	emailResource := models.MakeResource(recipient.ID, newEmail.Name)

	var thread *models.Thread
	err = env.Threads.WhereAndFetchOne(map[string]interface{}{
		"name":  rawSubject,
		"owner": recipient.ID,
	}, &thread)
	if err != nil {
		thread = &models.Thread{
			Resource: models.MakeResource(recipient.ID, rawSubject),
			Emails:   []string{emailResource.ID},
			Labels:   []string{inbox.ID},
			Members:  append(append(newEmail.To, newEmail.CC...), newEmail.BCC...),
			IsRead:   false,
		}

		err := env.Threads.Insert(thread)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to create a new thread")
			return
		}
	} else {
		existingMembers := make(map[string]struct{})

		for _, member := range thread.Members {
			existingMembers[member] = struct{}{}
		}

		for _, member := range newEmail.To {
			if _, ok := existingMembers[member]; !ok {
				thread.Members = append(thread.Members, member)
				existingMembers[member] = struct{}{}
			}
		}

		for _, member := range newEmail.CC {
			if _, ok := existingMembers[member]; !ok {
				thread.Members = append(thread.Members, member)
				existingMembers[member] = struct{}{}
			}
		}

		for _, member := range newEmail.BCC {
			if _, ok := existingMembers[member]; !ok {
				thread.Members = append(thread.Members, member)
				existingMembers[member] = struct{}{}
			}
		}

		err := env.Threads.UpdateID(thread.ID, thread)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    thread.ID,
				"error": err.Error(),
			}).Error("Unable to update an existing thread")
			return
		}
	}

	// Insert the new email
	newEmail.Resource = emailResource
	newEmail.Status = "processed"
	newEmail.Thread = thread.ID

	err = env.Emails.Insert(newEmail)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to create a new email")
		return
	}

	// Send notifications
	err = env.Producer.Publish("email_delivery", map[string]interface{}{
		"id":    email.ID,
		"owner": email.Owner,
	})
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"id":    email.ID,
			"error": err.Error(),
		}).Error("Unable to publish a delivery message")
	}

	err = env.NATS.Publish("email_receipt", map[string]interface{}{
		"id":    newEmail.ID,
		"owner": newEmail.Owner,
	})
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"id":    newEmail.ID,
			"error": err.Error(),
		}).Error("Unable to publish a receipt message")
	}
}*/
