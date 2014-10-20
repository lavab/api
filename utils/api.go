package utils

import (
	"encoding/json"
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
	fmt.Fprint(w, packIt(data))
}

//InternalErrorResponse TODO
func ErrorResponse(w http.ResponseWriter, code int, message string, debug string) {
	out := map[string]interface{}{
		"status":  code,
		"message": message,
		"debug":   debug,
		"success": false,
	}
	if debug == "" {
		delete(out, "debug")
	}
	fmt.Fprint(w, packIt(out))
}

// packIt receives a generic map and tries to marshal it
func packIt(data map[string]interface{}) string {
	res, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error marshalling the response body.", err)
		return "{\"status\": 500, \"message\":\"Error while marshalling the response body\"}"
	}
	return string(res)
}
