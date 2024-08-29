package app

import (
	"log/slog"

	"github.com/blankspace9/notes-app/internal/app/httpapp"
	"github.com/blankspace9/notes-app/internal/config"
	"github.com/blankspace9/notes-app/internal/delivery/rest"
	"github.com/blankspace9/notes-app/internal/external/spellchecker"
	"github.com/blankspace9/notes-app/internal/services/authservice"
	noteservice "github.com/blankspace9/notes-app/internal/services/noteService"
	"github.com/blankspace9/notes-app/internal/storage"
)

type App struct {
	HTTPServer *httpapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	storage, err := storage.New(storage.PostgresConnectionInfo(cfg.Storage))
	if err != nil {
		panic(err)
	}

	authService := authservice.New(log, storage, storage, cfg.JWT)

	spellChecker := spellchecker.New(cfg.SpellChecker.URL)
	notesService := noteservice.New(log, storage, spellChecker)

	handler := rest.New(log, authService, notesService)

	httpApp := httpapp.New(log, handler.InitRouter(), cfg.HTTPServer.Port, cfg.HTTPServer.Timeout)

	return &App{
		HTTPServer: httpApp,
	}
}
