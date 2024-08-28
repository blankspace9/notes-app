package app

import (
	"log/slog"

	httpapp "github.com/blankspace9/notes-app/internal/app/http"
	"github.com/blankspace9/notes-app/internal/config"
)

type App struct {
	HTTPServer *httpapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	// TODO: storage init
	// TODO: services init
	// TODO: handler init
	// TODO: app init

	return &App{}
}
