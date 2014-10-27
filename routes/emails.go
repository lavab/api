package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lavab/api/utils"
)

// Emails TODO
func Emails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// CreateEmail TODO
func CreateEmail(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.ReadJSON(r.Body)
	if err != nil {
		utils.ErrorResponse(w, 400, "Couldn't parse the request body", err.Error())
	}
	log.Println(reqData)

}

// Email TODO
func Email(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// UpdateEmail TODO
func UpdateEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// DeleteEmail TODO
func DeleteEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
