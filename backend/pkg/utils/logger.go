package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global zerolog configuration
func InitLogger(env string) {
	zerolog.TimeFieldFormat = time.RFC3339

	if env == "development" || env == "local" {
		// Use console writer for human-readable output in local/dev
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// Use JSON format for production
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msgf("Logger initialized in %s mode", env)
}
