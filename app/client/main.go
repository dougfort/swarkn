package main

import (
	"os"

	"github.com/rs/zerolog"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	logger := zerolog.New(os.Stdout).
		Level(zerolog.DebugLevel).
		With().Timestamp().Str("service", "swarkn").Logger()

	if err := run(logger); err != nil {
		logger.Error().AnErr("main", err).Msg("exit with error")
		os.Exit(1)
	}
}

func run(logger zerolog.Logger) error {
	logger.Info().Msg("program starts)")

	defer logger.Info().Msg("program exits)")

	return nil
}
