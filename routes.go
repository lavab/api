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
	route{"/logout", routes.Logout, "DELETE", ""},
	route{"/me", routes.Me, "GET", ""},
	route{"/me", routes.UpdateMe, "PUT", ""},
	route{"/me/wipe-user-data", routes.WipeUserData, "DELETE", ""},
	route{"/me/delete-account", routes.DeleteAccount, "DELETE", ""},
	route{"/threads", routes.Threads, "GET", ""},
	route{"/threads/{id}", routes.Thread, "GET", ""},
	route{"/threads/{id}", routes.UpdateThread, "PUT", ""},
	route{"/emails", routes.Emails, "GET", ""},
	route{"/emails", routes.CreateEmail, "POST", ""},
	route{"/emails/{id}", routes.Email, "GET", ""},
	route{"/emails/{id}", routes.UpdateEmail, "PUT", ""},
	route{"/emails/{id}", routes.DeleteEmail, "DELETE", ""},
	route{"/labels", routes.Labels, "GET", ""},
	route{"/labels", routes.CreateLabel, "POST", ""},
	route{"/labels/{id}", routes.Label, "GET", ""},
	route{"/labels/{id}", routes.UpdateLabel, "PUT", ""},
	route{"/labels/{id}", routes.DeleteLabel, "DELETE", ""},
	route{"/contacts", routes.Contacts, "GET", ""},
	route{"/contacts", routes.CreateContact, "POST", ""},
	route{"/contacts/{id}", routes.Contact, "GET", ""},
	route{"/contacts/{id}", routes.UpdateContact, "PUT", ""},
	route{"/contacts/{id}", routes.DeleteContact, "DELETE", ""},
	route{"/keys", routes.SubmitKey, "POST", ""},
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
