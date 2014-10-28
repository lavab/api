package utils

import (
	"fmt"
	"net/http"
)

// JSONResponse writes JSON to an http.ResponseWriter with the corresponding status code
func JSONResponse(w http.ResponseWriter, status int, data map[string]interface{}) {
	if status < 100 || status > 599 {
		status = 200
	}
	w.WriteHeader(status)
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
