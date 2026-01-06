package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/piotmni/go-mini-templates/ent-atlas/ent"
	"github.com/piotmni/go-mini-templates/ent-atlas/internal/config"
	"github.com/piotmni/go-mini-templates/ent-atlas/internal/db"
	internalhttp "github.com/piotmni/go-mini-templates/ent-atlas/internal/http"
	"github.com/piotmni/go-mini-templates/ent-atlas/internal/service"
)

type App struct {
	cfg    *config.Config
	orm    *ent.Client
	server *internalhttp.Server
}

func New(ctx context.Context) (*App, error) {
	// init zerolog logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// load config
	log.Debug().Msg("initializing app")
	log.Debug().Msg("loading config")
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// initialize ent ORM
	log.Debug().Msg("initializing ORM")
	ormClient, err := db.NewPostgres(cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	// Initialize services
	pasteService := service.NewPasteService(ormClient)

	// Initialize router
	router := internalhttp.NewRouter(internalhttp.RouterDeps{
		PasteService: pasteService,
	})

	// Initialize HTTP server
	server := internalhttp.NewServer(
		":"+cfg.Server.Port,
		router,
		15*time.Second,
		15*time.Second,
	)

	return &App{
		cfg:    cfg,
		orm:    ormClient,
		server: server,
	}, nil
}

// Run starts the application and blocks until shutdown.
func (a *App) Run() error {
	// Channel to listen for errors from server
	errChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Info().Str("port", a.cfg.Server.Port).Msg("starting HTTP server")
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
		log.Info().Str("signal", sig.String()).Msg("shutting down")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	if err := a.orm.Close(); err != nil {
		log.Error().Err(err).Msg("database close error")
	}

	log.Info().Msg("shutdown complete")
	return nil
}
