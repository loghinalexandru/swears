package handler

import (
	"encoding/json"
	"net/http"

	"github.com/loghinalexandru/swears/internal/encoding"
	"github.com/loghinalexandru/swears/internal/service"
	"github.com/rs/zerolog"
)

type Response struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

type RandomHandler struct {
	logger zerolog.Logger
	swears *service.Swears
}

func NewRandom(logger zerolog.Logger, svc *service.Swears) *RandomHandler {
	return &RandomHandler{
		logger: logger,
		swears: svc,
	}
}

func (handler *RandomHandler) Random(writer http.ResponseWriter, request *http.Request) {
	lang := "en"

	if request.URL.Query().Has("lang") {
		lang = request.URL.Query().Get("lang")
	}

	swear, err := handler.swears.GetSwear(lang)

	if err != nil {
		handler.logger.Err(err).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(Response{
		Swear: swear,
		Lang:  lang,
	})

	if err != nil {
		handler.logger.Err(err).Send()
	}

	writer.Write(res)
}

func (handler *RandomHandler) RandomFile(writer http.ResponseWriter, request *http.Request) {
	lang := "en"
	var enc service.Encoder

	if request.URL.Query().Has("lang") {
		lang = request.URL.Query().Get("lang")
	}

	switch request.URL.Query().Get("encoder") {
	case "opus":
		enc = encoding.NewOpus()
	}

	result := handler.swears.GetSwearFile(lang, enc)

	if result == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write(result)
}
