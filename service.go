package main

import (
	"os"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

type SwearsRepo interface {
	Get() string
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

	return repo.Get()
}

func (svc SwearsSvc) GetSwearFile(lang string) []byte {
	var result []byte
	repo, exits := svc.data[lang]

	if !exits {
		return result
	}

	// TODO: Check if generated already and send. Do not delete at the end
	config := htgotts.Speech{Folder: "misc", Language: voices.Romanian}
	config.CreateSpeechFile(repo.Get(), "temp")

	result, err := os.ReadFile("misc/temp.mp3")

	if err != nil {
		panic("Could not read file!")
	}

	defer os.Remove("misc/temp.mp3")
	return result
}
