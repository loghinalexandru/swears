package models

import (
	"github.com/google/uuid"
)

type Record struct {
	ID    uuid.UUID
	Value string
}
