package servehttp

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func Serve(ctx context.Context, logger zerolog.Logger, serverErrors chan<- error) {
	const serverName = "http"
	// TODO: get port from env/cli
	const port = 3000
	const hello = "Hello World!"

	logger = logger.With().Str("server", serverName).Logger()

	handler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger.Debug().Str("method", r.Method).Str("URI", r.RequestURI).Msg("request")

		if strings.ToUpper(r.Method) == "POST" {
			// Handle POST by returning the hash of the Body contents to the caller
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
		} else {
			// Handle anything other than POST like the example service: send hello
			w.Header().Add("content-type", "text/plain")
			if _, err := w.Write([]byte(hello)); err != nil {
				logger.Error().AnErr("w.Write", err).Msg("writing hello to reeponse body")
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
