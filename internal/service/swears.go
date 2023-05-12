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

	"github.com/loghinalexandru/swears/internal/codec"
	"github.com/loghinalexandru/swears/internal/model"
)

var (
	errMissingRepo = errors.New("missing repository for asked language")
)

const (
	ttsURL = "http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s"
)

type swearsOpt func(*Swears)

type Swears struct {
	downloadPath string
	client       *http.Client
	mtx          sync.Mutex
	data         map[string]model.SwearsRepo
}

func NewSwears(repos []model.SwearsRepo, downloadPath string, opts ...swearsOpt) *Swears {
	result := &Swears{
		downloadPath: downloadPath,
		client:       http.DefaultClient,
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

func (svc *Swears) GetSwearFile(lang string, codec codec.Encoder) ([]byte, error) {
	repo, exists := svc.data[lang]

	if !exists {
		return nil, errMissingRepo
	}

	swear, err := repo.Get()

	if err != nil {
		return nil, err
	}

	fname := fmt.Sprintf("%s/%s.mp3", svc.downloadPath, swear.ID)
	_, err = os.Stat(fname)

	if os.IsNotExist(err) {
		svc.downloadTTSFile(fname, swear.Value, lang)
	}

	result, err := os.ReadFile(fname)

	if codec != nil {
		result, err = codec.Encode(bytes.NewReader(result))
	}

	if err != nil {
		return nil, err
	}

	return result, nil
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
