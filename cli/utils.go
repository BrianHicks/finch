package main

import (
	"log"
	"path/filepath"

	"github.com/brianhicks/finch"
	"github.com/syndtr/goleveldb/leveldb/storage"

	"os"
)

// getStorage takes the FINCH_STORAGE environment variable into account. If
// it's "mem", this will return an in-memory database. If that's not true,
// it'll return an instantiated TaskDB instance.
func getTaskDB() (*finch.TaskDB, error) {
	dbPath := os.Getenv("FINCH_STORAGE")
	if dbPath == "mem" {
		return finch.NewTaskDB(storage.NewMemStorage())
	}

	if dbPath == "" {
		home := os.Getenv("HOME")
		dbPath = filepath.Join(home, ".finchdb")
		log.Print(dbPath)
	}

	store, err := storage.OpenFile(dbPath)
	if err != nil {
		return new(finch.TaskDB), err
	}

	tdb, err := finch.NewTaskDB(store)
	if err != nil {
		return tdb, err
	}

	return tdb, nil
}
