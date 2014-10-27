package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/lavab/api/models"
)

// CurrentSession returns the current request's session object
func CurrentSession(r *http.Request) *models.Session {
	session, ok := context.Get(r, "session").(*models.Session)
	if !ok {
		log.Fatalln("Session data in gorilla/context was not found or malformed.")
	}
	return session
}

// JSONResponse packs a generic map and writes it to a http.ResponseWriter as JSON.
// Additionally, if there's a data["status"] value, it's going to be added as a header.
func JSONResponse(w http.ResponseWriter, data map[string]interface{}) {
	if code, ok := data["status"]; ok {
		w.WriteHeader(code.(int))
	}
	fmt.Fprint(w, MakeJSON(data))
}

// ErrorResponse TODO
func ErrorResponse(w http.ResponseWriter, code int, message string, debug string) {
	out := map[string]interface{}{
		"debug":   debug,
		"message": message,
		"status":  code,
		"success": false,
	}
	if debug == "" {
		delete(out, "debug")
	}
	fmt.Fprint(w, MakeJSON(out))
}

// FormatNotRecognizedResponse TODO
func FormatNotRecognizedResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, 400, "Format not recognized", err.Error())
}
