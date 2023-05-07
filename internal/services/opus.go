package services

import (
	"io"

	"github.com/jonas747/dca"
)

var (
	stdOptions = dca.StdEncodeOptions
)

func init() {
	stdOptions.RawOutput = true
}

func Encode(stream io.Reader) ([]byte, error) {
	encodeSession, err := dca.EncodeMem(stream, stdOptions)

	if err != nil {
		return nil, err
	}

	defer encodeSession.Cleanup()
	return io.ReadAll(encodeSession)
}
