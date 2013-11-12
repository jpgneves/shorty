package storage

import (
	"errors"
	"fmt"
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
	t.Flush()
}

func (t *TextDB) Flush() {
	f, err := os.Open(t.filename)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range t.data {
		fmt.Fprintf(f, "%v,%v\n", k, v)
	}
	f.Sync()
	f.Close()
}

func OpenDB(backend, datasource string) (DB, error) {
	if backend != "plain" {
		return nil, errors.New("Unsupported backend")
	} else if len(datasource) <= 0 {
		return nil, errors.New("Bad file")
	}
	return new(TextDB), nil
}