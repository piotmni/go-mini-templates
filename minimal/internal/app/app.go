package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/piotmni/go-mini-templates/minimal/internal/config"
	"github.com/piotmni/go-mini-templates/minimal/internal/db"
	"github.com/piotmni/go-mini-templates/minimal/internal/http"
	"github.com/piotmni/go-mini-templates/minimal/internal/http/handlers"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/note"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *pgx.Conn
	server *http.Server
}

func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	ctx := context.Background()

	// Initialize database
	database, err := db.NewPostgres(ctx, cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	categoryRepo := category.NewPostgresRepository(database)
	noteRepo := note.NewPostgresRepository(database)

	// Initialize services
	categoryService := category.NewService(categoryRepo, logger)
	noteService := note.NewService(noteRepo, logger)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	noteHandler := handlers.NewNoteHandler(noteService)

	// Initialize HTTP server
	server := http.NewServer(
		http.ServerConfig{
			Host:  cfg.Server.Host,
			Port:  cfg.Server.Port,
			Token: cfg.Auth.Token,
		},
		logger,
		categoryHandler,
		noteHandler,
	)

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     database,
		server: server,
	}, nil
}

func (a *App) Run() error {
	// Channel for server errors
	errChan := make(chan error, 1)

	// Start HTTP server in a goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		a.logger.Info("received shutdown signal", zap.String("signal", sig.String()))
	}

	// Graceful shutdown
	return a.Shutdown()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("failed to shutdown HTTP server", zap.Error(err))
	}

	// Close database connection
	if err := a.db.Close(ctx); err != nil {
		a.logger.Error("failed to close database connection", zap.Error(err))
	}

	a.logger.Info("application stopped")
	return nil
}
