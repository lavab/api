package utils

import (
	"net/http"

	"github.com/lavab/api/env"
)

// marshalJSON is a wrapper over json.Marshal that returns an error message if an error occured
func marshalJSON(data interface{}) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return `{"status":500,"message":"Error occured while marshalling the response body"}`, err
	}
	return result, nil
}

// JSONResponse writes JSON to an http.ResponseWriter with the corresponding status code
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	// Get rid of the invalid status codes
	if status < 100 || status > 599 {
		status = 200
	}

	// Try to marshal the input
	result, err := marshalJSON(data)
	if err != nil {
		env.G.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to marshal a message")
	}

	// Write the result
	w.WriteHeader(status)
	w.Write(result)
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
	w.Write(marshalJSON(out))
}

// FormatNotRecognizedResponse TODO
func FormatNotRecognizedResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, 400, "Format not recognized", err.Error())
}
