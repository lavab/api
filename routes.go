package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lavab/api/routes"
)

type handleFunc func(w http.ResponseWriter, r *http.Request)

type route struct {
	Path        string     `json:"path"`
	HandleFunc  handleFunc `json:"-"`
	Method      string     `json:"method"`
	Description string     `json:"desc"`
}

var publicRoutes = []route{
	route{"/", listRoutes, "GET", "List of public and auth methods"},
	route{"/keys/{id}", routes.Key, "GET", ""},
	route{"/login", routes.Login, "POST", ""},
	route{"/signup", routes.Signup, "POST", ""},
}

var authRoutes = []route{
	route{"/logout", routes.Logout, "DELETE", "Destroys the current session"},
	route{"/me", routes.Me, "GET", "Fetch profile data for the current user"},
	route{"/me", routes.UpdateMe, "PUT", "Update data for the current user (settings, billing data, password, etc.)"},
	route{"/me/sessions", routes.Sessions, "GET", "Lists all the active sessions for the current user"},
	route{"/me/wipe-user-data", routes.WipeUserData, "DELETE", "Deletes all personal data of the user, except for basic profile information and billing status"},
	route{"/me/delete-account", routes.DeleteAccount, "DELETE", "Permanently deletes the user account"},
	route{"/threads", routes.Threads, "GET", "List email threads for the current user"},
	route{"/threads/{id}", routes.Thread, "GET", "Fetch a specific email thread"},
	route{"/threads/{id}", routes.UpdateThread, "PUT", "Update an email thread"},
	route{"/emails", routes.Emails, "GET", "List all emails for the current user"},
	route{"/emails", routes.CreateEmail, "POST", "Create and send an email"},
	route{"/emails/{id}", routes.Email, "GET", "Fetch a specific email"},
	route{"/emails/{id}", routes.UpdateEmail, "PUT", "Update a specific email (label, archive, etc)"},
	route{"/emails/{id}", routes.DeleteEmail, "DELETE", "Delete an email"},
	route{"/labels", routes.Labels, "GET", "List labels for the current user"},
	route{"/labels", routes.CreateLabel, "POST", "Create a new label"},
	route{"/labels/{id}", routes.Label, "GET", "Fetch a specific label"},
	route{"/labels/{id}", routes.UpdateLabel, "PUT", "Update a label"},
	route{"/labels/{id}", routes.DeleteLabel, "DELETE", "Delete a label"},
	route{"/contacts", routes.Contacts, "GET", "List all contacts for the current user"},
	route{"/contacts", routes.CreateContact, "POST", "Create a new contact"},
	route{"/contacts/{id}", routes.Contact, "GET", "Fetch a specific contact"},
	route{"/contacts/{id}", routes.UpdateContact, "PUT", "Update a contact"},
	route{"/contacts/{id}", routes.DeleteContact, "DELETE", "Delete a contact"},
	route{"/keys", routes.SubmitKey, "POST", "Submit a key to the Lavaboom private server"},
	route{"/keys/{id}", routes.VoteKey, "POST", "Vote or flag a key"},
}

func listRoutes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, config.RootJSON)
}

func rootResponseString() string {
	tmp, err := json.Marshal(
		map[string]interface{}{
			"message":  "Lavaboom API",
			"docs_url": "http://lavaboom.readme.io/",
			"version":  cApiVersion,
			"routes": map[string]interface{}{
				"public": publicRoutes,
				"auth":   authRoutes,
			}})
	if err != nil {
		log.Fatalln("Error! Couldn't marshal JSON.", err)
	}
	return string(tmp)
}
