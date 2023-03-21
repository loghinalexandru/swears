package main

import (
	"log"
	"net/http"
	"os"

	"github.com/loghinalexandru/swears/internal/handlers"
	"github.com/loghinalexandru/swears/internal/models"
	"github.com/loghinalexandru/swears/internal/repository"
	"github.com/loghinalexandru/swears/internal/services"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	roRepo := repository.New(logger, "ro", "misc/datastore/ro.txt")
	frRepo := repository.New(logger, "fr", "misc/datastore/fr.txt")
	enRepo := repository.New(logger, "en", "misc/datastore/en.txt")

	svc := services.NewSwears(
		[]models.SwearsRepo{
			roRepo,
			frRepo,
			enRepo,
		},
		http.DefaultClient,
		"misc")

	handler := handlers.NewRandom(logger, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/random", tracing(logger, contentType(handler.Random, "application/json")))
	mux.HandleFunc("/api/random/file", tracing(logger, contentType(handler.RandomFile, "application/octet-stream")))

	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	server.ListenAndServe()
}

func contentType(next http.HandlerFunc, mediaType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", mediaType)
		next(w, r)
	}
}

func tracing(logger *log.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Calling endpoint %v", r.URL)
		next(w, r)
	}
}
