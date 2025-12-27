package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

// AuthMiddleware validates HTTP Basic Auth with username + API key
func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}

		// Decode base64 username:apikey
		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}
		username, apiKey := parts[0], parts[1]

		// Validate username + API key
		ctx := r.Context()
		user, err := h.usersService.GetUserByUsername(ctx, username)
		if err != nil || user == nil || user.APIKey != apiKey {
			http.Error(w, "invalid username or API key", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "username", user.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
