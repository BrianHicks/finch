package persist

import (
	"bytes"

	"github.com/syndtr/goleveldb/leveldb"
)

// Range implements a range over some LevelDB keys. Since this always has to
// know what DB to range over, use DB.Key() and similar to construct.
type Range struct {
	Start []byte
	Limit []byte
	store *Store
}

func (r *Range) contains(target []byte) bool {
	return bytes.Compare(r.Start, target) <= 0 &&
		bytes.Compare(r.Limit, target) > 0
}

// First gets the first matching value in the Range
func (r *Range) First() ([]byte, error) {
	iter := r.store.DB.NewIterator(r.store.RO)
	defer iter.Release()

	iter.Seek(r.Start)

	if !r.contains(iter.Key()) {
		return []byte{}, leveldb.ErrNotFound
	}

	return iter.Value(), nil
}

// Last gets the last matching value in the Range
func (r *Range) Last() ([]byte, error) {
	iter := r.store.DB.NewIterator(r.store.RO)
	defer iter.Release()

	iter.Seek(r.Limit)
	iter.Prev()

	if !r.contains(iter.Key()) {
		return []byte{}, leveldb.ErrNotFound
	}

	return iter.Value(), nil
}

// All returns a slice of all the values in the range.
func (r *Range) All() ([][]byte, error) {
	iter := r.store.DB.NewIterator(r.store.RO)
	defer iter.Release()

	values := [][]byte{}
	for ok := iter.Seek(r.Start); ok; ok = iter.Next() {
		if !r.contains(iter.Key()) {
			break
		}
		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())
		values = append(values, value)
	}

	return values, iter.Error()
}
