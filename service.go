package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	tts "github.com/hegedustibor/htgo-tts"
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
		config := tts.Speech{Folder: "misc", Language: lang}
		config.CreateSpeechFile(repo.Get().Value, strconv.Itoa(repo.Get().Index))
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
