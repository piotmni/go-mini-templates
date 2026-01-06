package main

import (
	"log"

	"github.com/piotmni/go-mini-templates/minimal/internal/app"
	"github.com/piotmni/go-mini-templates/minimal/internal/config"
	"github.com/piotmni/go-mini-templates/minimal/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLogger, err := logger.New("info", true)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer zapLogger.Sync()

	application, err := app.New(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("failed to create application", zap.Error(err))
	}

	if err := application.Run(); err != nil {
		zapLogger.Fatal("application error", zap.Error(err))
	}
}
