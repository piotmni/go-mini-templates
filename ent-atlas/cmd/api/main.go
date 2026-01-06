package main

import (
	"context"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/piotmni/go-mini-templates/ent-atlas/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize application")
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("application error")
		os.Exit(1)
	}
}
