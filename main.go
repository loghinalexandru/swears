package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/loghinalexandru/swears/repository"
)

type handler func(w http.ResponseWriter, r *http.Request)

type Response struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func randomHandler(svc SwearsSvc) handler {
	return contentTypeJSON(func(w http.ResponseWriter, r *http.Request) {
		lang := "en"

		if r.URL.Query().Has("lang") {
			lang = r.URL.Query().Get("lang")
		}

		swear := svc.GetSwear(lang)

		if swear == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res, err := json.Marshal(Response{
			Swear: swear,
			Lang:  lang,
		})

		if err != nil {
			log.Fatal(err)
		}

		w.Write(res)
	})
}

func soundFileHandler(svc SwearsSvc) handler {
	return contentTypeMP3(func(w http.ResponseWriter, r *http.Request) {
		lang := "en"

		if r.URL.Query().Has("lang") {
			lang = r.URL.Query().Get("lang")
		}

		result := svc.GetSwearFile(lang)
		w.Write(result)
	})
}

func main() {
	roRepo := repository.New("ro", "misc/ro.txt")

	svc := NewSwears([]SwearsRepo{
		roRepo,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/random", randomHandler(svc))
	mux.HandleFunc("/api/random/file", soundFileHandler(svc))

	http.ListenAndServe(":3000", mux)
}

func contentTypeJSON(next handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next(w, r)
	}
}

func contentTypeMP3(next handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/octet-stream")
		next(w, r)
	}
}
