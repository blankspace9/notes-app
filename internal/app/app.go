package app

import (
	"log/slog"

	"github.com/blankspace9/notes-app/internal/app/httpapp"
	"github.com/blankspace9/notes-app/internal/config"
	"github.com/blankspace9/notes-app/internal/delivery/rest"
	"github.com/blankspace9/notes-app/internal/services/authservice"
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

	authService := authservice.New(log, storage, storage, cfg.Tokens, cfg.JWT_SECRET)

	handler := rest.New(log, authService, nil) // TODO: add notes service

	httpApp := httpapp.New(log, handler.InitRouter(), cfg.HTTPServer.Port, cfg.HTTPServer.Timeout)

	return &App{
		HTTPServer: httpApp,
	}
}
