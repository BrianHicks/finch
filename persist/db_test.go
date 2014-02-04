package persist

import (
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestFileLifecycle(t *testing.T) {
	name := "_taskdb_lifecycle"
	store, err := NewFile(name)
	assert.Nil(t, err)
	defer os.RemoveAll(name)

	store.Close()
}

// TestMemory just to make sure we don't cause panics or anything. It shouldn't
// error, otherwise.
func TestMemory(t *testing.T) {
	t.Parallel()
	_, err := NewInMemory()
	assert.Nil(t, err)
}

func TestRange(t *testing.T) {
	t.Parallel()
	store, err := NewInMemory()
	assert.Nil(t, err)

	r := store.Range([]byte("start"), []byte("end"))

	assert.Equal(t, r.Start, []byte("start"))
	assert.Equal(t, r.Limit, []byte("end"))
	assert.Equal(t, r.store, store)
}

func TestPrefixRange(t *testing.T) {
	t.Parallel()
	store, err := NewInMemory()
	assert.Nil(t, err)

	r := store.Prefix([]byte{0, 1})

	assert.Equal(t, r.Start, []byte{0, 1})
	assert.Equal(t, r.Limit, []byte{0, 2})
	assert.Equal(t, r.store, store)
}
