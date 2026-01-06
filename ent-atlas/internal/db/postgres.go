package db

import (
	"github.com/piotmni/go-mini-templates/ent-atlas/ent"
	"github.com/rs/zerolog/log"
)

func NewPostgres(connectionURL string) (*ent.Client, error) {

	log.Debug().Msg("initializing database connection")
	client, err := ent.Open("postgres", connectionURL)

	if err != nil {
		return nil, err
	}

	return client, nil
}
