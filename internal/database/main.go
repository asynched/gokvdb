package database

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Record struct {
	Value string `json:"value"`
	Ttl   int64  `json:"ttl"`
}

type Database struct {
	lock *sync.RWMutex
	data map[string]Record
}

func New() *Database {
	return &Database{
		lock: &sync.RWMutex{},
		data: make(map[string]Record),
	}
}

func (db *Database) Set(key string, value Record) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.data[key] = value
}

func (db *Database) Get(key string) (Record, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	value, ok := db.data[key]

	if !ok {
		return Record{}, false
	}

	if value.Ttl == 0 {
		return value, true
	}

	if value.Ttl < time.Now().Unix() {
		return Record{}, false
	}

	return value, true
}

func (db *Database) Delete(key string) {
	db.lock.Lock()
	defer db.lock.Unlock()

	delete(db.data, key)
}

func (db *Database) FlushAll() {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.data = make(map[string]Record)
}

func (db *Database) Snapshot() map[string]Record {
	db.lock.RLock()
	defer db.lock.RUnlock()

	data := make(map[string]Record)

	for k, v := range db.data {
		data[k] = v
	}

	return data
}

func (db *Database) Load(filename string) error {
	defer db.lock.Unlock()
	db.lock.Lock()

	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	data := make(map[string]Record)
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return err
	}

	db.data = data

	return nil
}

func (db *Database) Dump(filename string) error {
	defer db.lock.RUnlock()
	db.lock.RLock()

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(db.data)

	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
