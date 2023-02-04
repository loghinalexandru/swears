package repository

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"
)

type fileDB struct {
	lang string
	data []string
}

func New(language string, path string) *fileDB {
	db := &fileDB{
		lang: language,
		data: []string{},
	}
	db.Load(path)

	return db
}

func (db *fileDB) Lang() string {
	return db.lang
}

func (db *fileDB) Get() string {
	gen := rand.New(rand.NewSource(time.Now().UnixMilli()))
	index := gen.Intn(len(db.data))

	return db.data[index]
}

func (db *fileDB) Load(filePath string) {
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
		db.data = append(db.data, scaner.Text())
	}
}
