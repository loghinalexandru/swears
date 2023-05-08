package main

import (
	"net/http"
	"os"

	"github.com/loghinalexandru/swears/internal/handler"
	"github.com/loghinalexandru/swears/internal/model"
	"github.com/loghinalexandru/swears/internal/repository"
	"github.com/loghinalexandru/swears/internal/service"
	"github.com/rs/zerolog"
)

const (
	storagePath = "misc"
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

	svc := service.NewSwears(
		[]model.SwearsRepo{
			roRepo,
			frRepo,
			enRepo,
		},
		storagePath,
		service.WithLogger(logger),
	)

	handlerRand := handler.NewRandom(logger, svc)
	handlerRemote := handler.NewRemote(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/random", logRoute(logger, contentType(handlerRand.Random, "application/json")))
	mux.HandleFunc("/api/random/file", logRoute(logger, contentType(handlerRand.RandomFile, "application/octet-stream")))
	mux.HandleFunc("/api/remote", logRoute(logger, contentType(handlerRemote.RemoteVideo, "application/octet-stream")))

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
			Str("method", r.Method).
			Str("addr", r.RemoteAddr).
			Str("agent", r.UserAgent()).
			Str("url", r.URL.Path).
			Msg("Request received")
		next(w, r)
	}
}
