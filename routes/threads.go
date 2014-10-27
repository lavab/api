package routes

import (
	"fmt"
	"net/http"
)

// Threads TODO
func Threads(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// Thread TODO
func Thread(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// UpdateThread TODO
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
