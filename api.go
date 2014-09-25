package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lavab/api/routes"
)

func init() {
}

func main() {
	r := mux.NewRouter()

	if host := os.Getenv("API_HOST"); host != "" {
		r = r.Schemes("https").Host(host).Subrouter()
	}

	r.HandleFunc("/", routes.Root).Methods("GET")
	r.HandleFunc("/login", routes.Login).Methods("POST")
	r.HandleFunc("/signup", routes.Signup).Methods("POST")
	r.HandleFunc("/logout", routes.Logout)
	r.HandleFunc("/me", routes.Me).Methods("GET", "PUT")

	r.HandleFunc("/actions/wipe-user-data", routes.WipeUserData).Methods("DELETE")
	r.HandleFunc("/actions/delete-account", routes.DeleteAccount).Methods("DELETE")

	r.HandleFunc("/threads", routes.Threads).Methods("GET", "POST")
	r.HandleFunc("/threads/{id}", routes.Thread).Methods("GET", "PUT")

	r.HandleFunc("/messages", routes.Messages).Methods("GET", "POST")
	r.HandleFunc("/messages/{id}", routes.Message).Methods("GET", "DELETE", "PUT")

	r.HandleFunc("/tags", routes.Tags).Methods("GET", "POST")
	r.HandleFunc("/tags/{id}", routes.Tag).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/tags/{id}/threads", routes.TagThreads).Methods("GET")
	r.HandleFunc("/tags/{id}/messages", routes.TagMessages).Methods("GET")

	r.HandleFunc("/contacts", routes.Contacts).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", routes.Contact).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/contacts/{id}/threads", routes.ContactThreads).Methods("GET")

	r.HandleFunc("/keys", routes.Keys).Methods("GET")
	r.HandleFunc("/keys/{id}", routes.Key).Methods("GET")
	r.HandleFunc("/keys/{id}/jwt", routes.KeyJwt).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
