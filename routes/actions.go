package routes

import (
	"fmt"
	"net/http"
)

// WipeUserData TODO
func WipeUserData(w http.ResponseWriter, r *http.Request) {

}

// DeleteAccount TODO
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
