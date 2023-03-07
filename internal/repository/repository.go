package repository

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/loghinalexandru/swears/internal/models"
)

const (
	emptyDataStore = "empty data store"
)

type fileDB struct {
	lang   string
	logger *log.Logger
	data   []models.Record
}

func New(logger *log.Logger, language string, path string) *fileDB {
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
		return models.Record{}, errors.New(emptyDataStore)
	}

	index := rand.Intn(len(db.data))
	return db.data[index], nil
}

func (db *fileDB) load(filePath string) {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		panic("No file found!")
	}

	fh, err := os.Open(filePath)
	if err != nil {
		db.logger.Fatal("Could not open file!")
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
