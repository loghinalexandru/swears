package repository

import (
	"bufio"
	"errors"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/loghinalexandru/swears/internal/models"
	"github.com/rs/zerolog"
)

var (
	ErrEmptyDataStore = errors.New("empty data store")
)

type fileDB struct {
	lang   string
	logger zerolog.Logger
	data   []models.Record
}

func New(logger zerolog.Logger, language string, path string) *fileDB {
	db := &fileDB{
		lang:   language,
		logger: logger,
		data:   []models.Record{},
	}
	db.load(path)

	return db
}

func (db *fileDB) Lang() string {
	return db.lang
}

func (db *fileDB) Get() (models.Record, error) {
	if len(db.data) == 0 {
		return models.Record{}, ErrEmptyDataStore
	}

	index := rand.Intn(len(db.data))
	return db.data[index], nil
}

func (db *fileDB) load(filePath string) {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		db.logger.Panic().Msg("no file found")
	}

	fh, err := os.Open(filePath)
	if err != nil {
		db.logger.Panic().Msg("could not open file")
	}

	defer fh.Close()

	scaner := bufio.NewScanner(fh)
	scaner.Split(bufio.ScanLines)

	for scaner.Scan() {
		db.data = append(db.data, models.Record{
			ID:    uuid.New(),
			Value: scaner.Text(),
		})
	}
}
