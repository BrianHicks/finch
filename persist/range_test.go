package persist

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRangeContains(t *testing.T) {
	t.Parallel()

	r := Range{
		Start: []byte{1},
		Limit: []byte{3},
	}

	assert.True(t, r.Contains([]byte{1}))
	assert.True(t, r.Contains([]byte{2}))
	assert.False(t, r.Contains([]byte{3}))
}

func TestRangeFirst(t *testing.T) {
	t.Parallel()
	db, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte("test")
	doc := []byte{1}

	err = db.Put(key, doc, db.wo)
	assert.Nil(t, err)

	// get the value back out?
	ret, err := db.Prefix(key).First()
	assert.Nil(t, err)

	assert.Equal(t, doc, ret)

	// and if there is no such value, error
	_, err = db.Prefix([]byte("asdf")).First()
	assert.Equal(t, ErrNoResult, err)
}

func TestRangeLast(t *testing.T) {
	t.Parallel()
	db, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte("test/2")
	doc := []byte{1}

	err = db.Put([]byte("test/1"), []byte{0}, db.wo)
	assert.Nil(t, err)
	err = db.Put(key, doc, db.wo)
	assert.Nil(t, err)

	// get the value back out?
	ret, err := db.Prefix(key).Last()
	assert.Nil(t, err)

	assert.Equal(t, doc, ret)

	// but still if there is no value, error
	_, err = db.Prefix([]byte("asdf")).Last()
	assert.Equal(t, ErrNoResult, err)
}

// There's a special case where calling Last should get the First value (if
// there's only one.) We'll put documents on either side to make sure it
// returns the right one.
func TestRangeLastIsFirst(t *testing.T) {
	t.Parallel()
	db, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte{2}
	doc := []byte{1}

	err = db.Put([]byte{1}, []byte{}, db.wo)
	assert.Nil(t, err)
	err = db.Put([]byte{3}, []byte{}, db.wo)
	assert.Nil(t, err)

	err = db.Put(key, doc, db.wo)
	assert.Nil(t, err)

	// get the value back out?
	lst, err := db.Prefix(key).Last()
	assert.Nil(t, err)

	fst, err := db.Prefix(key).First()
	assert.Nil(t, err)

	assert.Equal(t, doc, fst)
	assert.Equal(t, doc, lst)
}
