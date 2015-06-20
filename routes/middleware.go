package routes

import (
	"net/http"
	"strings"

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
			utils.JSONResponse(w, 401, utils.NewError(
				utils.MiddlewareMissingToken, "Missing auth token", false,
			))
			return
		}

		// Split it into two parts
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.JSONResponse(w, 401, utils.NewError(
				utils.MiddlewareInvalidFormat, "Invalid authorization header's format", false,
			))
			return
		}

		// Get the token from the database
		token, err := env.Tokens.GetToken(headerParts[1])
		if err != nil {
			utils.JSONResponse(w, 401, utils.NewError(
				utils.MiddlewareInvalidToken, err, false,
			))
			return
		}

		// Check if it's expired
		if token.Expired() {
			utils.JSONResponse(w, 401, utils.NewError(
				utils.MiddlewareExpiredToken, "Your authentication token has expired", false,
			))
			env.Tokens.DeleteID(token.ID)
			return
		}

		// Continue to the next middleware/route
		c.Env["token"] = token
		h.ServeHTTP(w, r)
	})
}
