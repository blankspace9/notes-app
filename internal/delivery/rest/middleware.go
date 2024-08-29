package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/auth"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
)

// Authentification
// Checks if the token is valid, and decides to let it go to endpoints or not
func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid token request: "+err.Error(), http.StatusUnauthorized)
			h.log.Warn("invalid token request", sl.Err(err))
			return
		}

		userId, expiresAt, err := h.authService.ParseToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Failed to parse token: "+err.Error(), http.StatusUnauthorized)
			h.log.Warn("failed to parse token", sl.Err(err))
			return
		}

		if expiresAt.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			h.log.Warn("token expired")
			return
		}

		ctx := context.WithValue(r.Context(), auth.CtxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Retrieving a token from a request
func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty Authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
