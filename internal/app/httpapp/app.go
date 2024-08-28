package httpapp

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/blankspace9/notes-app/pkg/httpserver"
)

type App struct {
	log        *slog.Logger
	httpServer httpserver.Server
	port       int
}

func New(log *slog.Logger, httpHandler http.Handler, port int, timeout time.Duration) *App {
	httpServer := httpserver.NewServer(httpHandler, httpserver.Port(fmt.Sprint(port)), httpserver.ReadTimeout(timeout), httpserver.WriteTimeout(timeout), httpserver.ShutdownTimeout(timeout))

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

	log.Info("HTTP server is running", slog.Int("port", a.port))
}

func (a *App) Stop() {
	const op = "httpapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping HTTP server", slog.Int("port", a.port))

	a.httpServer.Shutdown()
}

func (a *App) Notify() <-chan error {
	return a.httpServer.Notify()
}
