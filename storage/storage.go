package storage

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
)

// Because I can't compile sqlite3 right now, let's just do a dumb text dump
type DB interface {
	Find(key string) interface{}
	Insert(key string, value interface{})
	Flush()
}

type TextDB struct {
	filename string
	data map[string]interface{}
}

func (t *TextDB) Find(key string) interface{} {
	if v, ok := t.data[key]; ok {
		return v
	}
	return nil
}

func (t *TextDB) Insert(key string, value interface{}) {
	t.data[key] = value
}

func (t *TextDB) Flush() {
	f, err := os.OpenFile(t.filename, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	for k, v := range t.data {
		line := []string{k, v.(string)}
		log.Printf("%v storing %v = %v as %v", t.filename, k, v, line)
		err = writer.Write(line)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer writer.Flush()
	f.Sync()
}

func OpenDB(backend, datasource string) (DB, error) {
	if backend != "plain" {
		return nil, errors.New("Unsupported backend")
	} else if len(datasource) <= 0 {
		return nil, errors.New("Bad file")
	}
	db := &TextDB{datasource, make(map[string]interface{})}
	f, err := os.OpenFile(datasource, os.O_RDONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	entries, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for idx := range entries {
		if len(entries[idx]) >= 2 {
			db.Insert(entries[idx][0], entries[idx][1])
		}
	}
	return db, nil
}