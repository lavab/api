package utils

import (
	"fmt"
	"net/http"
)

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
