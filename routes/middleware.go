package routes

import (
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/utils"
)

// AuthMiddlewareResponse is the response sent by the middleware if user is not logged in
type AuthMiddlewareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AuthMiddleware checks whether the token passed with the request is valid
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

		// Get the token from the database
		token, err := env.TokensCache.GetToken(headerParts[1])
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Cannot retrieve session from the database")

			utils.JSONResponse(w, 401, &AuthMiddlewareResponse{
				Success: false,
				Message: "Non existing authorization token",
			})
			return
		}

		// We don't need to invalidate the those in db since
		// Redis expires them automatically, probably need a
		// goroutine that checks for the expired ones ?

		// Continue to the next middleware/route
		c.Env["session"] = token
		h.ServeHTTP(w, r)
	})
}
