package main

import (
	"path/filepath"

	"github.com/BrianHicks/finch/core"
	"github.com/syndtr/goleveldb/leveldb/storage"

	"os"
)

// getStorage takes the FINCH_STORAGE environment variable into account. If
// it's "mem", this will return an in-memory database. If that's not true,
// it'll return an instantiated TaskDB instance.
func getTaskDB() (*core.TaskDB, error) {
	dbPath := os.Getenv("FINCH_STORAGE")
	if dbPath == "mem" {
		return core.NewTaskDB(storage.NewMemStorage())
	}

	if dbPath == "" {
		home := os.Getenv("HOME")
		dbPath = filepath.Join(home, ".finchdb")
	}

	store, err := storage.OpenFile(dbPath)
	if err != nil {
		return new(core.TaskDB), err
	}

	tdb, err := core.NewTaskDB(store)
	if err != nil {
		return tdb, err
	}

	return tdb, nil
}
