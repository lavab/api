package routes

import (
	"fmt"
	"net/http"

	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// Emails TODO
func Emails(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
	mock := map[string]interface{}{
		"n_items": 1,
		"emails":  []models.Email{},
	}
	utils.JSONResponse(w, mock)
}

// CreateEmail TODO
func CreateEmail(w http.ResponseWriter, r *http.Request) {
	// reqData, err := utils.ReadJSON(r.Body)
	// if err != nil {
	// utils.ErrorResponse(w, 400, "Couldn't parse the request body", err.Error())
	// }
	// log.Println(reqData)
	mock := map[string]interface{}{
		"success": true,
		"created": []string{utils.UUID()},
	}
	utils.JSONResponse(w, mock)
}

// Email TODO
func Email(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
	mock := map[string]interface{}{
		"status": "sending",
	}
	utils.JSONResponse(w, mock)
}

// UpdateEmail TODO
func UpdateEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}

// DeleteEmail TODO
func DeleteEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"success\":false,\"message\":\"Sorry, not implemented yet\"}")
}
