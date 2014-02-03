package persist

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

// DB wraps a LevelDB instance and sets sane defaults
type DB struct {
	*leveldb.DB
	WO *opt.WriteOptions
	RO *opt.ReadOptions
}

// newDB takes a storage and returns DB instance
func New(store storage.Storage) (*DB, error) {
	db := new(DB)

	// Open the Database with the provided Storage
	options := &opt.Options{
		Filter: filter.NewBloomFilter(15),
	}
	DB, err := leveldb.Open(store, options)
	if err != nil {
		return db, err
	}
	db.DB = DB

	// Set default read and write options
	db.WO = &opt.WriteOptions{
		Sync: true,
	}
	db.RO = &opt.ReadOptions{
		DontFillCache: false,
	}

	return db, nil
}

func NewFile(fname string) (*DB, error) {
	store, err := storage.OpenFile(fname)
	if err != nil {
		return new(DB), err
	}
	return New(store)
}

func NewInMemory() (*DB, error) {
	return New(storage.NewMemStorage())
}

func (db *DB) Range(start, end []byte) *Range {
	return &Range{
		Start: start,
		Limit: end,
		db:    db,
	}
}

func (db *DB) Prefix(start []byte) *Range {
	end := make([]byte, len(start))
	size := copy(end, start)

	end[size-1] = end[size-1] + 1

	return &Range{
		Start: start,
		Limit: end,
		db:    db,
	}
}
