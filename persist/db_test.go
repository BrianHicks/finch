package persist

import (
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestFileLifecycle(t *testing.T) {
	name := "_taskdb_lifecycle"
	db, err := NewFile(name)
	assert.Nil(t, err)
	defer os.RemoveAll(name)

	db.Close()
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
	db, err := NewInMemory()
	assert.Nil(t, err)

	r := db.Range([]byte("start"), []byte("end"))

	assert.Equal(t, r.Start, []byte("start"))
	assert.Equal(t, r.End, []byte("end"))
	assert.Equal(t, r.db, db)
}

func TestPrefixRange(t *testing.T) {
	t.Parallel()
	db, err := NewInMemory()
	assert.Nil(t, err)

	r := db.Prefix([]byte{0, 1})

	assert.Equal(t, r.Start, []byte{0, 1})
	assert.Equal(t, r.End, []byte{0, 2})
	assert.Equal(t, r.db, db)
}
