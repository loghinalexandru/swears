package main

import (
	"net/http"
)

func randomHandler() func(w http.ResponseWriter, r *http.Request) {
	svc := NewSwears()
	return func(w http.ResponseWriter, r *http.Request) {
		lang := "ro"

		if r.URL.Query().Has("lang") {
			lang = r.URL.Query().Get("lang")
		}

		w.Write([]byte(svc.GetSwear(lang)))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/random", randomHandler())

	http.ListenAndServe(":3333", mux)
}
