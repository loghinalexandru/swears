package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jonas747/dca"
	"github.com/loghinalexandru/swears/models"
)

type SwearsRepo interface {
	Get() models.Record
	Lang() string
	Load(file string)
}

type SwearsSvc struct {
	data map[string]SwearsRepo
}

func NewSwears(repos []SwearsRepo) SwearsSvc {
	result := SwearsSvc{
		data: make(map[string]SwearsRepo),
	}

	for _, repo := range repos {
		result.data[repo.Lang()] = repo
	}

	return result
}

func (svc SwearsSvc) GetSwear(lang string) string {
	repo, exits := svc.data[lang]

	if !exits {
		return ""
	}

	return repo.Get().Value
}

func (svc SwearsSvc) GetSwearFile(lang string, opus bool) []byte {
	var result []byte
	repo, exits := svc.data[lang]

	if !exits {
		return result
	}

	fname := fmt.Sprintf("misc/%s.mp3", strconv.Itoa(repo.Get().Index))
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		downloadTTSFile(fname, repo.Get().Value, lang)
	}

	if opus {
		encdOpt := dca.StdEncodeOptions
		encdOpt.RawOutput = true
		encodeSession, err := dca.EncodeFile(fname, encdOpt)

		if err != nil {
			fmt.Println(err)
		}

		fname = fmt.Sprintf("misc/%s.dca", strconv.Itoa(repo.Get().Index))
		output, err := os.Create(fname)
		if err != nil {
			panic(err)
		}

		io.Copy(output, encodeSession)
		output.Close()
		encodeSession.Cleanup()
	}

	result, err = os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	return result
}

func downloadTTSFile(fileName string, text string, lang string) error {
	f, err := os.Open(fileName)
	if err != nil {
		url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), lang)
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		output, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(output, response.Body)
		return err
	}

	f.Close()
	return nil
}
