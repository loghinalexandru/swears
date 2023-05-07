package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/loghinalexandru/swears/internal/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	errVideoToLong = errors.New("video is too long to play")
)

const (
	durationMax = time.Second * 30
)

type RemoteHandler struct {
	logger zerolog.Logger
	client youtube.Client
}

func NewRemote(logger zerolog.Logger) *RemoteHandler {
	return &RemoteHandler{
		logger: logger,
		client: youtube.Client{},
	}
}

func (h *RemoteHandler) RemoteVideo(writer http.ResponseWriter, request *http.Request) {
	var ID string
	var opus bool
	var result []byte

	if request.URL.Query().Has("id") {
		ID = request.URL.Query().Get("id")
	}

	if request.URL.Query().Has("opus") {
		opus, _ = strconv.ParseBool(request.URL.Query().Get("opus"))
	}

	metadata, err := h.client.GetVideo(ID)

	if err != nil {
		log.Err(err).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if metadata.Duration > durationMax {
		log.Err(errVideoToLong).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	formats := metadata.Formats.WithAudioChannels()
	stream, _, err := h.client.GetStream(metadata, &formats[0])

	if err != nil {
		log.Err(err).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer stream.Close()
	result, err = io.ReadAll(stream)

	if opus {
		result, err = services.Encode(bytes.NewReader(result))
	}

	if err != nil {
		log.Err(err).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(result)
}
