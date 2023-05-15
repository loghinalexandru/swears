package model

import (
	"github.com/google/uuid"
)

type Record struct {
	ID    uuid.UUID
	Value string
}

type SwearsRepo interface {
	Get() (Record, error)
	Lang() string
}
