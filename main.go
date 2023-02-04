package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request)

type Response struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func randomHandler(svc SwearsSvc) handler {
	return contentType(func(w http.ResponseWriter, r *http.Request) {
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

func main() {
	svc := NewSwears()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/random", randomHandler(svc))

	http.ListenAndServe(":3000", mux)
}

func contentType(next handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next(w, r)
	}
}
