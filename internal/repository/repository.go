package repository

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/loghinalexandru/swears/internal/model"
)

var (
	errEmptyDataStore = errors.New("empty data store")
)

type fileDB struct {
	lang string
	data []model.Record
}

func New(language string, path string) *fileDB {
	db := &fileDB{
		lang: language,
		data: []model.Record{},
	}
	db.load(path)

	return db
}

func (db *fileDB) Lang() string {
	return db.lang
}

func (db *fileDB) Get() (model.Record, error) {
	if len(db.data) == 0 {
		return model.Record{}, errEmptyDataStore
	}

	index := rand.Intn(len(db.data))
	return db.data[index], nil
}

func (db *fileDB) load(filePath string) {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("no file found for path: %q", filePath))
	}

	fh, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("could not open file for path: %q", filePath))
	}

	defer fh.Close()

	scaner := bufio.NewScanner(fh)
	scaner.Split(bufio.ScanLines)

	for scaner.Scan() {
		db.data = append(db.data, model.Record{
			ID:    uuid.New(),
			Value: scaner.Text(),
		})
	}
}
