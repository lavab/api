package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

import "crypto/rand"

// JSONResponse just marshals a map[string]interface{} to json
// Additionally, if the map contains a "status" key, it adds an HTTP error code
func JSONResponse(w http.ResponseWriter, data map[string]interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Woot! Error marshaling the response body", 500)
		return
	}
	if code, ok := data["status"]; ok {
		w.WriteHeader(code.(int))
	}
	fmt.Fprintf(w, string(res))
}

// HoursFromNow returns time.Now + n hours
// It uses RFC3339 and UTC to yield comparable strings
// Example: 2006-01-02T15:04:05Z00:00
func HoursFromNow(n int) string {
	return time.Now().UTC().Add(time.Hour * time.Duration(n)).Format(time.RFC3339)
}

// RandomString returns a secure random string of a certain length
func RandomString(length int) (string, error) {
	tmp := make([]byte, length)
	_, err := rand.Read(tmp)
	if err != nil {
		return "", err
	}
	return string(tmp), nil
}

// FileExists is a stupid little wrapper of os.Stat that checks whether a file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}
