package routes

import (
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"

	"github.com/lavab/api/db"
	"github.com/lavab/api/dbutils"
	"github.com/lavab/api/utils"
)

type AuthMiddlewareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func AuthMiddleware(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the Authorization header
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.JSONResponse(w, 401, &AuthMiddlewareResponse{
				Success: false,
				Message: "Missing auth token",
			})
			return
		}

		// Split it into two parts
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.JSONResponse(w, 401, &AuthMiddlewareResponse{
				Success: false,
				Message: "Invalid authorization header",
			})
			return
		}

		// Get the session from the database
		session, ok := dbutils.GetSession(headerParts[1])
		if !ok {
			utils.JSONResponse(w, 401, &AuthMiddlewareResponse{
				Success: false,
				Message: "Invalid authorization token",
			})
			return
		}

		// Check if it's expired
		if session.Expired() {
			utils.JSONResponse(w, 419, &AuthMiddlewareResponse{
				Success: false,
				Message: "Authorization token has expired",
			})
			db.Delete("sessions", session.ID)
			return
		}

		// Continue to the next middleware/route
		c.Env["session"] = session
		h.ServeHTTP(w, r)
	})
}
