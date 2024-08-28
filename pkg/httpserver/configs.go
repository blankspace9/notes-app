package httpserver

import (
	"net"
	"time"
)

type Config func(*Server)

func Port(port string) Config {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Config {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Config {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Config {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
