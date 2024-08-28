package app

import (
	"log/slog"

	httpapp "github.com/blankspace9/notes-app/internal/app/http"
	"github.com/blankspace9/notes-app/internal/config"
	"github.com/blankspace9/notes-app/internal/delivery/rest"
)

type App struct {
	HTTPServer *httpapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	// TODO: storage init
	// TODO: services init

	handler := rest.New(log, nil, nil) // TODO: add services

	httpApp := httpapp.New(log, handler.InitRouter(), cfg.HTTPServer.Port, cfg.HTTPServer.Timeout)

	return &App{
		HTTPServer: httpApp,
	}
}
