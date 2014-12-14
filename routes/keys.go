package routes

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// KeysListResponse contains the result of the KeysList request
type KeysListResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Keys    *[]string `json:"keys,omitempty"`
}

// KeysList responds with the list of keys assigned to the spiecified email
func KeysList(w http.ResponseWriter, r *http.Request) {
	// Get the username from the GET query
	user := r.URL.Query().Get("user")
	if user == "" {
		utils.JSONResponse(w, 409, &KeysListResponse{
			Success: false,
			Message: "Invalid username",
		})
		return
	}

	account, err := env.Accounts.FindAccountByName(user)
	if err != nil {
		utils.JSONResponse(w, 409, &KeysListResponse{
			Success: false,
			Message: "Invalid username",
		})
		return
	}

	// Find all keys owner by user
	keys, err := env.Keys.FindByOwner(account.ID)
	if err != nil {
		utils.JSONResponse(w, 500, &KeysListResponse{
			Success: false,
			Message: "Internal server error (KE/LI/01)",
		})
		return
	}

	// Equivalent of _.keys(keys) in JavaScript with underscore.js
	keyIDs := []string{}
	for _, key := range keys {
		keyIDs = append(keyIDs, key.ID)
	}

	// Respond with list of keys
	utils.JSONResponse(w, 200, &KeysListResponse{
		Success: true,
		Keys:    &keyIDs,
	})
}

// KeysCreateRequest contains the data passed to the KeysCreate endpoint.
type KeysCreateRequest struct {
	Key string `json:"key" schema:"key"` // gpg armored key
}

// KeysCreateResponse contains the result of the KeysCreate request.
type KeysCreateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Key     *models.Key `json:"key,omitempty"`
}

// KeysCreate appens a new key to the server
func KeysCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input KeysCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 409, &KeysCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Get the session
	session := c.Env["token"].(*models.Token)

	// Parse the armored key
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(input.Key))
	if err != nil {
		utils.JSONResponse(w, 409, &KeysCreateResponse{
			Success: false,
			Message: "Invalid key format",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
			"list":  entityList,
		}).Warn("Cannot parse an armored key")
		return
	}

	// Parse using armor pkg
	block, err := armor.Decode(strings.NewReader(input.Key))
	if err != nil {
		utils.JSONResponse(w, 409, &KeysCreateResponse{
			Success: false,
			Message: "Invalid key format",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
			"list":  entityList,
		}).Warn("Cannot parse an armored key #2")
		return
	}

	// Get the account from db
	account, err := env.Accounts.GetAccount(session.Owner)
	if err != nil {
		utils.JSONResponse(w, 500, &KeysCreateResponse{
			Success: false,
			Message: "Internal server error - KE/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
			"id":    session.Owner,
		}).Error("Cannot fetch user from database")
		return
	}

	// Let's hope that the user is capable of sending proper armored keys
	publicKey := entityList[0]

	// Encode the fingerprint
	id := hex.EncodeToString(publicKey.PrimaryKey.Fingerprint[:])

	// Get the key's bit length - should not return an error
	bitLength, _ := publicKey.PrimaryKey.BitLength()

	// Allocate a new key
	key := &models.Key{
		Resource: models.MakeResource(
			account.ID,
			fmt.Sprintf(
				"%s/%d/%s",
				utils.GetAlgorithmName(publicKey.PrimaryKey.PubKeyAlgo),
				bitLength,
				publicKey.PrimaryKey.KeyIdString(),
			),
		),
		Headers:     block.Header,
		Algorithm:   utils.GetAlgorithmName(publicKey.PrimaryKey.PubKeyAlgo),
		Length:      bitLength,
		Key:         input.Key,
		KeyID:       publicKey.PrimaryKey.KeyIdString(),
		KeyIDShort:  publicKey.PrimaryKey.KeyIdShortString(),
		Reliability: 0,
	}

	// Update id as we can't do it directly during allocation
	key.ID = id

	// Try to insert it into the database
	if err := env.Keys.Insert(key); err != nil {
		utils.JSONResponse(w, 500, &KeysCreateResponse{
			Success: false,
			Message: "Internal server error - KE/CR/02",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not insert a key to the database")
		return
	}

	// Return the inserted key
	utils.JSONResponse(w, 201, &KeysCreateResponse{
		Success: true,
		Message: "A new key has been successfully inserted",
		Key:     key,
	})
}

// KeysGetResponse contains the result of the KeysGet request.
type KeysGetResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Key     *models.Key `json:"key,omitempty"`
}

// KeysGet does *something* - TODO
func KeysGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get ID from the passed URL params
	id := c.URLParams["id"]

	// Fetch the requested key from the database
	key, err := env.Keys.FindByFingerprint(id)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to fetch the requested key from the database")

		utils.JSONResponse(w, 404, &KeysGetResponse{
			Success: false,
			Message: "Requested key does not exist on our server",
		})
		return
	}

	// Return the requested key
	utils.JSONResponse(w, 200, &KeysGetResponse{
		Success: true,
		Key:     key,
	})
}

// KeysVoteResponse contains the result of the KeysVote request.
type KeysVoteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// KeysVote does *something* - TODO
func KeysVote(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 501, &KeysVoteResponse{
		Success: false,
		Message: "Sorry, not implemented yet",
	})
}
