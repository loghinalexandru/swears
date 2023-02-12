package repository

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/loghinalexandru/swears/models"
)

const (
	emptyDataStore = "Empty data store!"
)

type fileDB struct {
	lang      string
	generator *rand.Rand
	data      []models.Record
}

func New(language string, path string) *fileDB {
	db := &fileDB{
		lang:      language,
		generator: rand.New(rand.NewSource(time.Now().UnixNano())),
		data:      []models.Record{},
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

	index := db.generator.Intn(len(db.data))
	return db.data[index], nil
}

func (db *fileDB) load(filePath string) {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		panic("No file found!")
	}

	fh, err := os.Open(filePath)
	defer fh.Close()

	if err != nil {
		log.Fatal("Could not open file!")
	}

	scaner := bufio.NewScanner(fh)
	scaner.Split(bufio.ScanLines)

	for scaner.Scan() {
		db.data = append(db.data, models.Record{
			ID:    uuid.New(),
			Value: scaner.Text(),
		})
	}
}
