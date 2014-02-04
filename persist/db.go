package persist

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

// DB wraps a LevelDB instance and sets sane defaults
type Store struct {
	DB *leveldb.DB
	WO *opt.WriteOptions
	RO *opt.ReadOptions
}

// newDB takes a storage and returns DB instance
func New(storage storage.Storage) (*Store, error) {
	store := new(Store)

	// Open the Database with the provided Storage
	options := &opt.Options{
		Filter: filter.NewBloomFilter(15),
	}
	DB, err := leveldb.Open(storage, options)
	if err != nil {
		return store, err
	}
	store.DB = DB

	// Set default read and write options
	store.WO = &opt.WriteOptions{
		Sync: true,
	}
	store.RO = &opt.ReadOptions{
		DontFillCache: false,
	}

	return store, nil
}

func NewFile(fname string) (*Store, error) {
	storage, err := storage.OpenFile(fname)
	if err != nil {
		return new(Store), err
	}
	return New(storage)
}

func NewInMemory() (*Store, error) {
	return New(storage.NewMemStorage())
}

func (store *Store) Close() {
	store.DB.Close()
	store.DB = nil
}

func (store *Store) Range(start, end []byte) *Range {
	return &Range{
		Start: start,
		Limit: end,
		store: store,
	}
}

func (store *Store) Prefix(start []byte) *Range {
	end := make([]byte, len(start))
	size := copy(end, start)

	end[size-1] = end[size-1] + 1

	return &Range{
		Start: start,
		Limit: end,
		store: store,
	}
}
