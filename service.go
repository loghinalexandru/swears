package main

import "github.com/loghinalexandru/swears/repository"

type SwearsRepo interface {
	Get() string
	Lang() string
	Load(file string)
}

type SwearsSvc struct {
	data map[string]SwearsRepo
}

func NewSwears() SwearsSvc {
	result := SwearsSvc{
		data: make(map[string]SwearsRepo),
	}

	repo := repository.New("ro", "misc/ro.txt")
	result.data[repo.Lang()] = repo

	return result
}

func (svc SwearsSvc) GetSwear(lang string) string {
	repo, exits := svc.data[lang]

	if !exits {
		return ""
	}

	return repo.Get()
}
