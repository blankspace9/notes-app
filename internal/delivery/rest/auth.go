package rest

import (
	"encoding/json"
	"net/http"

	"github.com/blankspace9/notes-app/internal/domain/auth"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
	"github.com/blankspace9/notes-app/internal/services/authservice"
)

func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	var authData auth.AuthInput

	d := json.NewDecoder(r.Body)
	err := d.Decode(&authData)
	if err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to parse json", sl.Err(err))
		return
	}

	// Validate fields
	err = authData.Validate()
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		h.log.Warn("invalid credentials", sl.Err(err))
		return
	}

	id, err := h.authService.RegisterUser(r.Context(), authData.Email, authData.Password)
	if err != nil {
		switch err {
		case authservice.ErrUserExists:
			http.Error(w, "Failed to register user: "+authservice.ErrUserExists.Error(), http.StatusConflict)
			h.log.Warn("failed to register user", sl.Err(err))
		default:
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			h.log.Warn("failed to register user", sl.Err(err))
		}
		return
	}

	resp, err := json.Marshal(map[string]int64{
		"id": id,
	})
	if err != nil {
		http.Error(w, "Failed to marshal response json: "+err.Error(), http.StatusInternalServerError)
		h.log.Warn("failed to marshal response json", sl.Err(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var authData auth.AuthInput

	d := json.NewDecoder(r.Body)
	err := d.Decode(&authData)
	if err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to parse json", sl.Err(err))
		return
	}

	// Validate fields
	err = authData.Validate()
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		h.log.Warn("invalid credentials", sl.Err(err))
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), authData.Email, authData.Password)
	if err != nil {
		switch err {
		case authservice.ErrInvalidCredentials:
			http.Error(w, "Failed to login: "+authservice.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
			h.log.Warn("failed to login", sl.Err(err))
		default:
			http.Error(w, "Failed to login: "+err.Error(), http.StatusInternalServerError)
			h.log.Warn("failed to login", sl.Err(err))
		}
		return
	}

	resp, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		http.Error(w, "Failed to marshal response json: "+err.Error(), http.StatusInternalServerError)
		h.log.Warn("failed to marshal response json", sl.Err(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		Path:     "/auth",
	})
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		h.log.Warn("invalid refresh token", sl.Err(err))
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Failed to refresh tokens", http.StatusBadRequest)
		h.log.Warn("failed to refresh tokens", sl.Err(err))
		return
	}

	// Form the answer to the client
	resp, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		http.Error(w, "Failed to marshal response json: "+err.Error(), http.StatusInternalServerError)
		h.log.Warn("failed to marshal response json", sl.Err(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		Path:     "/auth",
	})
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}
