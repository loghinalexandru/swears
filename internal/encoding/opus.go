package encoding

import (
	"io"

	"github.com/jonas747/dca"
)

type opus struct {
	options *dca.EncodeOptions
}

func NewOpus() *opus {
	opt := dca.StdEncodeOptions
	opt.RawOutput = true

	return &opus{
		options: opt,
	}
}

func (enc opus) Encode(stream io.Reader) ([]byte, error) {
	encodeSession, err := dca.EncodeMem(stream, enc.options)

	if err != nil {
		return nil, err
	}

	defer encodeSession.Cleanup()
	return io.ReadAll(encodeSession)
}
