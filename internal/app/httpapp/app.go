package httpapp

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/blankspace9/notes-app/pkg/httpserver"
)

type App struct {
	log        *slog.Logger
	httpServer httpserver.Server
	port       string
}

func New(log *slog.Logger, httpHandler http.Handler, port string, timeout time.Duration) *App {
	httpServer := httpserver.NewServer(httpHandler, httpserver.Port(port), httpserver.ReadTimeout(timeout), httpserver.WriteTimeout(timeout), httpserver.ShutdownTimeout(timeout))

	return &App{
		log:        log,
		httpServer: *httpServer,
		port:       port,
	}
}

func (a *App) Run() {
	const op = "httpapp.Run"

	log := a.log.With(slog.String("op", op))

	a.httpServer.Start()

	log.Info("HTTP server is running", slog.String("port", a.port))
}

func (a *App) Stop() {
	const op = "httpapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping HTTP server", slog.String("port", a.port))

	a.httpServer.Shutdown()
}

func (a *App) Notify() <-chan error {
	return a.httpServer.Notify()
}
