package finch

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestKeySerialization(t *testing.T) {
	k := Key{"test", "2014", "abc"}

	assert.Equal(t, []byte("test/2014/abc"), k.Serialize())
}

func TestKeyDeserialization(t *testing.T) {
	k := Key{"test", "2014", "abc"}

	// test good
	k2, err := DeserializeKey([]byte("test/2014/abc"))
	assert.Nil(t, err)
	assert.Equal(t, &k, k2)

	// test bad
	k2, err = DeserializeKey([]byte(""))
	assert.Equal(t, ErrDeserializeKey, err)
	assert.Equal(t, new(Key), k2)

	// test weird
	k2, err = DeserializeKey([]byte("a/b/c/d"))
	assert.Nil(t, err)
	assert.Equal(t, &Key{"a", "b", "c/d"}, k2)
}

func TestKeyForTask(t *testing.T) {
	task := Task{Added: time.Now(), Id: "test"}

	key := KeyForTask("idx", &task)

	assert.Equal(t, "idx", key.Index)
	assert.Equal(t, task.Added.Format(time.RFC3339), key.Timestamp)
	assert.Equal(t, task.Id, key.Hash)
}
