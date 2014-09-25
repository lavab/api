package routes

import "net/http"

// CheckAuth TODO this method will check each request for a valid token
func CheckAuth(r *http.Request) bool {
	return true
}
