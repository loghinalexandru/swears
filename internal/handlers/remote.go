package handlers

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
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

// AddLater: rewrite this with a service & opus toggle
func (h *RemoteHandler) RemoteVideo(writer http.ResponseWriter, request *http.Request) {
	var ID string

	if request.URL.Query().Has("id") {
		ID = request.URL.Query().Get("id")
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

	if err != nil {
		log.Err(err).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	encdOpt := dca.StdEncodeOptions
	encdOpt.RawOutput = true
	encodeSession, err := dca.EncodeMem(stream, encdOpt)

	if err != nil {
		h.logger.Err(err).Send()
	}

	result, err := io.ReadAll(encodeSession)
	defer encodeSession.Cleanup()

	if err != nil {
		log.Err(err).Send()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(result)
}
