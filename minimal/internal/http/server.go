package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/piotmni/go-mini-templates/minimal/internal/http/handlers"
	"github.com/piotmni/go-mini-templates/minimal/internal/http/middleware"
	"go.uber.org/zap"
)

// ServerConfig holds the HTTP server configuration.
type ServerConfig struct {
	Host  string
	Port  int
	Token string
}

// Server is the HTTP server.
type Server struct {
	echo            *echo.Echo
	cfg             ServerConfig
	logger          *zap.Logger
	categoryHandler *handlers.CategoryHandler
	noteHandler     *handlers.NoteHandler
}

// NewServer creates a new HTTP server.
func NewServer(
	cfg ServerConfig,
	logger *zap.Logger,
	categoryHandler *handlers.CategoryHandler,
	noteHandler *handlers.NoteHandler,
) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		echo:            e,
		cfg:             cfg,
		logger:          logger.Named("http.server"),
		categoryHandler: categoryHandler,
		noteHandler:     noteHandler,
	}
}

// setupRoutes configures all routes.
func (s *Server) setupRoutes() {
	// Global middleware
	s.echo.Use(echomiddleware.Recover())
	s.echo.Use(echomiddleware.RequestID())
	s.echo.Use(s.requestLogger())

	// Health check (no auth required)
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes with auth
	api := s.echo.Group("/api/v1")
	api.Use(middleware.Auth(s.cfg.Token))

	// Register routes
	categories := api.Group("/categories")
	s.categoryHandler.RegisterRoutes(categories)

	notes := api.Group("/notes")
	s.noteHandler.RegisterRoutes(notes)
}

// requestLogger returns a middleware that logs requests.
func (s *Server) requestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			s.logger.Info("request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", res.Status),
				zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
			)

			return nil
		}
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.setupRoutes()

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	s.logger.Info("starting HTTP server", zap.String("addr", addr))

	return s.echo.Start(addr)
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")
	return s.echo.Shutdown(ctx)
}
