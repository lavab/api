package routes

import (
	"fmt"
	"net/http"

	"github.com/lavab/api/utils"
)

// Contacts TODO
func Contacts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// CreateContact TODO
func CreateContact(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.ReadJSON(r.Body)
	if err == nil {
		utils.FormatNotRecognizedResponse(w, err)
		return
	}

}

// Contact TODO
func Contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// UpdateContact TODO
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// DeleteContact TODO
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
