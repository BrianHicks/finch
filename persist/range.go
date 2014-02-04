package persist

import (
	"bytes"
	"errors"
)

var (
	ErrNoResult = errors.New("no such result")
)

// Range implements a range over some LevelDB keys. Since this always has to
// know what DB to range over, use DB.Key() and similar to construct.
type Range struct {
	Start []byte
	Limit []byte
	store *Store
}

func (r *Range) Contains(target []byte) bool {
	return bytes.Compare(r.Start, target) <= 0 &&
		bytes.Compare(r.Limit, target) > 0
}

func (r *Range) First() ([]byte, error) {
	iter := r.store.DB.NewIterator(r.store.RO)
	defer iter.Release()

	iter.Seek(r.Start)

	if !r.Contains(iter.Key()) {
		return []byte{}, ErrNoResult
	}

	return iter.Value(), nil
}

func (r *Range) Last() ([]byte, error) {
	iter := r.store.DB.NewIterator(r.store.RO)
	defer iter.Release()

	iter.Seek(r.Limit)
	iter.Prev()

	if !r.Contains(iter.Key()) {
		return []byte{}, ErrNoResult
	}

	return iter.Value(), nil
}

func (r *Range) All() ([][]byte, error) {
	iter := r.store.DB.NewIterator(r.store.RO)
	defer iter.Release()

	values := [][]byte{}
	for ok := iter.Seek(r.Start); ok; ok = iter.Next() {
		if !r.Contains(iter.Key()) {
			break
		}
		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())
		values = append(values, value)
	}

	return values, iter.Error()
}
