package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/loghinalexandru/swears/internal/services"
)

type Response struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

type HTTPHandler struct {
	logger *log.Logger
	swears *services.Swears
}

func NewRandom(logger *log.Logger, svc *services.Swears) *HTTPHandler {
	return &HTTPHandler{
		logger: logger,
		swears: svc,
	}
}

func (handler *HTTPHandler) Random(writer http.ResponseWriter, request *http.Request) {
	lang := "en"

	if request.URL.Query().Has("lang") {
		lang = request.URL.Query().Get("lang")
	}

	swear, err := handler.swears.GetSwear(lang)

	if err != nil {
		handler.logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(Response{
		Swear: swear,
		Lang:  lang,
	})

	if err != nil {
		handler.logger.Println(err)
	}

	writer.Write(res)
}

func (handler *HTTPHandler) RandomFile(writer http.ResponseWriter, request *http.Request) {
	lang := "en"
	encode := false

	if request.URL.Query().Has("lang") {
		lang = request.URL.Query().Get("lang")
	}

	if request.URL.Query().Has("opus") {
		encode, _ = strconv.ParseBool(request.URL.Query().Get("opus"))
	}

	result := handler.swears.GetSwearFile(lang, encode)

	if result == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write(result)
}
