package storage

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"sync"
)

type KeyValue struct {
	Key string
	Value interface{}
}

// Because I can't compile sqlite3 right now, let's just do a dumb text dump
type DB interface {
	Iterator() chan *KeyValue
	Find(key string) interface{}
	Insert(key string, value interface{})
	Flush()
}

type TextDB struct {
	filename string
	mutex *sync.RWMutex
	data map[string]interface{}
}

func (t *TextDB) Iterator() chan *KeyValue {
	c := make(chan *KeyValue)
	go func() {
		for k, v := range t.data {
			t.mutex.RLock()
			c <- &KeyValue{k, v}
			t.mutex.RUnlock()
		}
		close(c)
	}()
	return c
}

func (t *TextDB) Find(key string) interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if v, ok := t.data[key]; ok {
		return v
	}
	return nil
}

func (t *TextDB) Insert(key string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data[key] = value
}

func (t *TextDB) Flush() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	f, err := os.OpenFile(t.filename, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	entry_count := len(t.data)
	entries := make([][]string, entry_count)
	idx := 0
	for k, v := range t.data {
		entries[idx] = []string{k, v.(string)}
		idx++
	}
	err = writer.WriteAll(entries)
	if err != nil {
		log.Fatal(err)
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
	db := &TextDB{datasource, new(sync.RWMutex), make(map[string]interface{})}
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