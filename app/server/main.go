package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"

	"github.com/dougfort/swarkn/config"
	"github.com/dougfort/swarkn/servehttp"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	logger := zerolog.New(os.Stdout).
		With().Timestamp().Str("service", "swarkn").Logger()

	cfg, err := config.LoadServerConfig()
	if err != nil {
		logger.Error().AnErr("LoadServerConfig", err).Msg("error loading config")
		os.Exit(1)
	}

	if err := run(logger, cfg); err != nil {
		logger.Error().AnErr("main", err).Msg("exit with error")
		os.Exit(1)
	}
}

func run(logger zerolog.Logger, cfg config.ServerConfig) error {
	logLevel := zerolog.InfoLevel
	if cfg.LogLevel == "debug" {
		logLevel = zerolog.DebugLevel
	}
	logger = logger.Level(logLevel)

	logger.Info().Str("loglevel", cfg.LogLevel).Msg("program starts")
	defer logger.Info().Msg("program exits)")

	ctx, cancel := context.WithCancel(context.Background())

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	logger.Info().Msg("serving http")
	go servehttp.Serve(ctx, logger, cfg, serverErrors)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		cancel()
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logger.Info().Str("signal", fmt.Sprintf("%v", sig)).Msg("shutdown")
		cancel()
	}

	return nil
}
