package encoding

import (
	"errors"
	"io"

	"github.com/jonas747/dca"
)

var (
	errUnexpectedFailure = errors.New("something went wrong with the encoding process")
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
	result, err := io.ReadAll(encodeSession)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errUnexpectedFailure
	}

	return result, nil

}
