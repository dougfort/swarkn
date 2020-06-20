package servehttp

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

func Serve(ctx context.Context, logger zerolog.Logger, serverErrors chan<- error) {
	const serverName = "http"
	// TODO: get port from env/cli
	const port = 3000

	logger = logger.With().Str("server", serverName).Logger()

	handler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger.Debug().Str("URI", r.RequestURI).Msg("request")

		hash := sha256.New()

		_, err := io.Copy(hash, r.Body)
		if err != nil {
			logger.Error().AnErr("io.Copy", err).Msg("copying request body")
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("content-type", "application/octet-stream")
			if _, err := w.Write(hash.Sum(nil)); err != nil {
				logger.Error().AnErr("w.Write", err).Msg("writing hash to reeponse body")
			}
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	logger.Info().Msgf("listening on: %s", addr)
	server := http.Server{Addr: addr, Handler: http.HandlerFunc(handler)}

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	<-ctx.Done()
	logger.Debug().Msg("ctx.Done")
	server.Close()
}
