package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/gorilla/mux"
)

type Handler struct {
	log          *slog.Logger
	authService  AuthService
	notesService NotesService
}

type AuthService interface {
	RegisterUser(ctx context.Context, email, password string) (userID int64, err error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
}

type NotesService interface {
	CreateNote(ctx context.Context, userID int64, text string) (noteID int64, err error)
	GetNotes(ctx context.Context, userID int64) (notes []models.Note, err error)
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
