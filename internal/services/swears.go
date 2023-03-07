package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/jonas747/dca"
	"github.com/loghinalexandru/swears/internal/models"
)

const (
	missingRepo = "missing repository for asked language"
	ttsURL      = "http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s"
)

type Swears struct {
	client *http.Client
	mtx    sync.Mutex
	data   map[string]models.SwearsRepo
}

func NewSwears(repos []models.SwearsRepo, client *http.Client) *Swears {
	result := Swears{
		client: client,
		mtx:    sync.Mutex{},
		data:   make(map[string]models.SwearsRepo),
	}

	for _, repo := range repos {
		result.data[repo.Lang()] = repo
	}

	return &result
}

func (svc *Swears) GetSwear(lang string) (string, error) {
	repo, exists := svc.data[lang]

	if !exists {
		return "", errors.New(missingRepo)
	}

	res, err := repo.Get()

	if err != nil {
		return "", err
	}

	return res.Value, nil
}

func (svc *Swears) GetSwearFile(lang string, opus bool) []byte {
	var result []byte
	repo, exists := svc.data[lang]

	if !exists {
		return result
	}

	swear, err := repo.Get()

	if err != nil {
		log.Println(err)
		return result
	}

	fname := fmt.Sprintf("misc/%s.mp3", swear.ID)
	_, err = os.Stat(fname)

	if os.IsNotExist(err) {
		svc.downloadTTSFile(fname, swear.Value, lang)
	}

	if opus {
		encdOpt := dca.StdEncodeOptions
		encdOpt.RawOutput = true
		encodeSession, err := dca.EncodeFile(fname, encdOpt)

		if err != nil {
			log.Println(err)
			return nil
		}

		fname = fmt.Sprintf("misc/%s.dca", swear.ID)
		output, err := os.Create(fname)
		if err != nil {
			log.Println(err)
			return nil
		}

		io.Copy(output, encodeSession)
		output.Close()
		encodeSession.Cleanup()
	}

	result, err = os.ReadFile(fname)
	if err != nil {
		log.Println(err)
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
