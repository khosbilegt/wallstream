package api

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware validates API keys from the Authorization header
func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract API key from "Bearer <api_key>" or just use the header value directly
		var apiKey string
		if strings.HasPrefix(authHeader, "Bearer ") {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}
			apiKey = parts[1]
		} else {
			apiKey = authHeader
		}

		if apiKey == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Get user by API key
		ctx := context.Background()
		user, err := h.usersService.GetUserByAPIKey(ctx, apiKey)
		if err != nil {
			http.Error(w, "invalid API key", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx = context.WithValue(r.Context(), "user_id", user.ID)
		ctx = context.WithValue(ctx, "username", user.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
