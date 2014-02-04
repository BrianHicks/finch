package persist

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

// Store wraps a LevelDB instance and sets sane defaults
type Store struct {
	DB *leveldb.DB
	WO *opt.WriteOptions
	RO *opt.ReadOptions
}

// New takes a storage and returns DB instance
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

// NewFile creates a Store from the filename specified.
func NewFile(fname string) (*Store, error) {
	storage, err := storage.OpenFile(fname)
	if err != nil {
		return new(Store), err
	}
	return New(storage)
}

// NewInMemory creates a Store in memory
func NewInMemory() (*Store, error) {
	return New(storage.NewMemStorage())
}

// Close out the database after use. Don't try and use the Store after you call
// this method!
func (store *Store) Close() {
	store.DB.Close()
	store.DB = nil
}

// Write commits a batch operation to the underlying storage
func (store *Store) Write(batch *leveldb.Batch) error {
	return store.DB.Write(batch, store.WO)
}

// Put writes a single key/value pair to the underlying storage (it logs this
// operation)
func (store *Store) Put(key, value []byte) error {
	batch := LoggedBatch{new(leveldb.Batch)}
	batch.Put(key, value)
	return store.Write(batch.Batch)
}

// Delete removes a key/value pair from the underlying storage (it logs this
// operation as a batch, so you will not receive confirmation if the delete
// actually did anything.)
func (store *Store) Delete(key []byte) error {
	batch := LoggedBatch{new(leveldb.Batch)}
	batch.Delete(key)
	return store.Write(batch.Batch)
}

// Get retrieves a single key/value pair from the database and returns an error
// if it cannot find it.
func (store *Store) Get(key []byte) ([]byte, error) {
	return store.DB.Get(key, store.RO)
}

// Range takes a start (inclusive) and limit (exclusive) and returns an object
// you can get specific documents within the range with.
func (store *Store) Range(start, limit []byte) *Range {
	return &Range{
		Start: start,
		Limit: limit,
		store: store,
	}
}

// Prefix is a special case of Range for iterating over keys prefixed with a
// given value.
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
