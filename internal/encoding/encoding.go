package encoding

import (
	"io"
	"strings"
)

type Encoder interface {
	Encode(io.Reader) ([]byte, error)
}

func FromString(encoderType string) Encoder {
	switch strings.ToLower(encoderType) {
	case "opus":
		return NewOpus()
	default:
		return nil
	}
}
