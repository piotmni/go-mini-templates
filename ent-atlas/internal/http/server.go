package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server wraps the HTTP server.
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP server.
func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

// Start begins listening for requests.
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
