package main

import (
	"net/http"
	"os"

	"github.com/loghinalexandru/swears/internal/handlers"
	"github.com/loghinalexandru/swears/internal/models"
	"github.com/loghinalexandru/swears/internal/repository"
	"github.com/loghinalexandru/swears/internal/services"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

	roRepo := repository.New(logger, "ro", "misc/datastore/ro.txt")
	frRepo := repository.New(logger, "fr", "misc/datastore/fr.txt")
	enRepo := repository.New(logger, "en", "misc/datastore/en.txt")

	svc := services.NewSwears(
		[]models.SwearsRepo{
			roRepo,
			frRepo,
			enRepo,
		},
		"misc",
		services.WithLogger(logger),
	)

	handler := handlers.NewRandom(logger, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/random", logRoute(logger, contentType(handler.Random, "application/json")))
	mux.HandleFunc("/api/random/file", logRoute(logger, contentType(handler.RandomFile, "application/octet-stream")))

	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	logger.Info().Msgf("Server starting on address %v", server.Addr)
	server.ListenAndServe()
}

func contentType(next http.HandlerFunc, mediaType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", mediaType)
		next(w, r)
	}
}

func logRoute(logger zerolog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info().
			Str("Method", r.Method).
			Str("RemoteAddr", r.RemoteAddr).
			Str("URL", r.URL.Path).
			Msg("Request received")
		next(w, r)
	}
}
