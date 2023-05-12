package codec

import (
	"io"
	"strings"
)

type Encoder interface {
	Encode(io.Reader) ([]byte, error)
}

func New(codecType string) Encoder {
	switch strings.ToLower(codecType) {
	case "opus":
		return newOpus()
	default:
		return nil
	}
}
