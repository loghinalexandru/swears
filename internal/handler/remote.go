package handler

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	videoID := parseID(request.URL.Query())
	metadata, err := h.client.GetVideo(videoID)

	if err != nil {
		log.Err(err).Msg("Unexpected error when retrieving video metadata")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if metadata.Duration > durationMax {
		log.Warn().Msgf("Provided video is longer than expected max length: %v", durationMax)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	formats := metadata.Formats.WithAudioChannels().Quality("tiny")
	formats.Sort()
	stream, _, err := h.client.GetStream(metadata, &formats[0])

	if err != nil {
		log.Err(err).Msg("Unexpected error when retrieving video data")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer stream.Close()
	result, err := io.ReadAll(stream)

	enc := parseEncoder(request.URL.Query())

	if enc != nil {
		result, err = enc.Encode(bytes.NewReader(result))
	}

	if err != nil {
		log.Err(err).Msg("Unexpected error when encoding video data")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(result)
}
