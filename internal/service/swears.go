package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/loghinalexandru/swears/internal/model"
	"github.com/rs/zerolog"
)

var (
	errMissingRepo = errors.New("missing repository for asked language")
	errEmptyBuffer = errors.New("resulting buffer has no data")
)

const (
	ttsURL = "http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s"
)

type Encoder interface {
	Encode(io.Reader) ([]byte, error)
}

type swearsOpt func(*Swears)

type Swears struct {
	downloadPath string
	logger       zerolog.Logger
	client       *http.Client
	mtx          sync.Mutex
	data         map[string]model.SwearsRepo
}

func NewSwears(repos []model.SwearsRepo, downloadPath string, opts ...swearsOpt) *Swears {
	result := &Swears{
		downloadPath: downloadPath,
		client:       http.DefaultClient,
		logger:       zerolog.Nop().With().Logger(),
		mtx:          sync.Mutex{},
		data:         make(map[string]model.SwearsRepo),
	}

	for _, repo := range repos {
		result.data[repo.Lang()] = repo
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

func WithLogger(logger zerolog.Logger) swearsOpt {
	return func(s *Swears) {
		s.logger = logger
	}
}

func WithClient(client *http.Client) swearsOpt {
	return func(s *Swears) {
		s.client = client
	}
}

func (svc *Swears) GetSwear(lang string) (string, error) {
	repo, exists := svc.data[lang]

	if !exists {
		return "", errMissingRepo
	}

	res, err := repo.Get()

	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// Return err and handle it in handler
func (svc *Swears) GetSwearFile(lang string, enc Encoder) []byte {
	var result []byte
	repo, exists := svc.data[lang]

	if !exists {
		return result
	}

	swear, err := repo.Get()

	if err != nil {
		svc.logger.Err(err).Send()
		return result
	}

	fname := fmt.Sprintf("%s/%s.mp3", svc.downloadPath, swear.ID)
	_, err = os.Stat(fname)

	if os.IsNotExist(err) {
		svc.downloadTTSFile(fname, swear.Value, lang)
	}

	result, err = os.ReadFile(fname)

	if enc != nil {
		result, err = enc.Encode(bytes.NewReader(result))
	}

	if err != nil {
		svc.logger.Err(err).Send()
		return nil
	}

	if len(result) == 0 {
		svc.logger.Err(errEmptyBuffer).Send()
		return nil
	}

	return result
}

func (svc *Swears) downloadTTSFile(fileName string, text string, lang string) error {
	url := fmt.Sprintf(ttsURL, url.QueryEscape(text), lang)

	response, err := svc.client.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	svc.mtx.Lock()
	defer svc.mtx.Unlock()

	output, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer output.Close()
	_, err = io.Copy(output, response.Body)

	return err
}
