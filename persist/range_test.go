package persist

import (
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"

	"testing"
)

func TestRangeContains(t *testing.T) {
	t.Parallel()

	r := Range{
		Start: []byte{1},
		Limit: []byte{3},
	}

	assert.True(t, r.contains([]byte{1}))
	assert.True(t, r.contains([]byte{2}))
	assert.False(t, r.contains([]byte{3}))
}

func TestRangeFirst(t *testing.T) {
	t.Parallel()
	store, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte("test")
	doc := []byte{1}

	err = store.DB.Put(key, doc, store.WO)
	assert.Nil(t, err)

	// get the value back out?
	ret, err := store.Prefix(key).First()
	assert.Nil(t, err)

	assert.Equal(t, doc, ret)

	// and if there is no such value, error
	_, err = store.Prefix([]byte("asdf")).First()
	assert.Equal(t, leveldb.ErrNotFound, err)
}

func TestRangeLast(t *testing.T) {
	t.Parallel()
	store, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte("test/2")
	doc := []byte{1}

	err = store.DB.Put([]byte("test/1"), []byte{0}, store.WO)
	assert.Nil(t, err)
	err = store.DB.Put(key, doc, store.WO)
	assert.Nil(t, err)

	// get the value back out?
	ret, err := store.Prefix(key).Last()
	assert.Nil(t, err)

	assert.Equal(t, doc, ret)

	// but still if there is no value, error
	_, err = store.Prefix([]byte("asdf")).Last()
	assert.Equal(t, leveldb.ErrNotFound, err)
}

// There's a special case where calling Last should get the First value (if
// there's only one.) We'll put documents on either side to make sure it
// returns the right one.
func TestRangeLastIsFirst(t *testing.T) {
	t.Parallel()
	store, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte{2}
	doc := []byte{1}

	err = store.DB.Put([]byte{1}, []byte{}, store.WO)
	assert.Nil(t, err)
	err = store.DB.Put([]byte{3}, []byte{}, store.WO)
	assert.Nil(t, err)

	err = store.DB.Put(key, doc, store.WO)
	assert.Nil(t, err)

	// get the value back out?
	lst, err := store.Prefix(key).Last()
	assert.Nil(t, err)

	fst, err := store.Prefix(key).First()
	assert.Nil(t, err)

	assert.Equal(t, doc, fst)
	assert.Equal(t, doc, lst)
}

func TestRangeAll(t *testing.T) {
	t.Parallel()
	store, err := NewInMemory()
	assert.Nil(t, err)

	prefix := byte(1)
	values := [][]byte{}

	for i := byte(1); i < 5; i++ {
		value := []byte{i}
		err := store.DB.Put([]byte{prefix, i}, value, store.WO)
		assert.Nil(t, err)
		values = append(values, value)
	}

	all, err := store.Prefix([]byte{prefix}).All()

	assert.Nil(t, err)
	assert.Equal(t, values, all)
}
