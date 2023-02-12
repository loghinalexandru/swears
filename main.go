package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/loghinalexandru/swears/repository"
)

type Response struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func randomHandler(svc *SwearsSvc) http.HandlerFunc {
	return contentType(func(w http.ResponseWriter, r *http.Request) {
		lang := "en"

		if r.URL.Query().Has("lang") {
			lang = r.URL.Query().Get("lang")
		}

		swear, err := svc.GetSwear(lang)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(Response{
			Swear: swear,
			Lang:  lang,
		})

		if err != nil {
			log.Println(err)
		}

		w.Write(res)
	}, "application/json")
}

func soundFileHandler(svc *SwearsSvc) http.HandlerFunc {
	return contentType(func(w http.ResponseWriter, r *http.Request) {
		lang := "en"
		encode := false

		if r.URL.Query().Has("lang") {
			lang = r.URL.Query().Get("lang")
		}

		if r.URL.Query().Has("opus") {
			encode, _ = strconv.ParseBool(r.URL.Query().Get("opus"))
		}

		result := svc.GetSwearFile(lang, encode)

		if result == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write(result)
	}, "application/octet-stream")
}

func main() {
	roRepo := repository.New("ro", "misc/datastore/ro.txt")
	frRepo := repository.New("fr", "misc/datastore/fr.txt")
	enRepo := repository.New("en", "misc/datastore/en.txt")

	svc := NewSwearsSvc([]SwearsRepo{
		roRepo,
		frRepo,
		enRepo,
	}, http.DefaultClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/random", randomHandler(svc))
	mux.HandleFunc("/api/random/file", soundFileHandler(svc))

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
