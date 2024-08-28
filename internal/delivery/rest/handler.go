package rest

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	log          *slog.Logger
	authService  AuthService
	notesService NotesService
}

type AuthService interface {
}

type NotesService interface {
}

func New(log *slog.Logger, as AuthService, ns NotesService) *Handler {
	return &Handler{
		log:          log,
		authService:  as,
		notesService: ns,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	{
		auth := api.PathPrefix("/auth").Subrouter()
		{
			auth.HandleFunc("/registration", h.registration).Methods(http.MethodPost)
			auth.HandleFunc("/login", h.login).Methods(http.MethodPost)
		}

		notes := api.PathPrefix("/notes").Subrouter()
		{
			notes.HandleFunc("", h.addNote).Methods(http.MethodPost)
			notes.HandleFunc("", h.getNotes).Methods(http.MethodGet)
		}
	}

	return r
}
