package main

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
