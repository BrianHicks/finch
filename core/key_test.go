package core

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestKeySerialization(t *testing.T) {
	k := Key{"2014", "abc"}

	assert.Equal(t, []byte("test/2014/abc"), k.Serialize("test"))
}

func TestKeyDeserialization(t *testing.T) {
	k := Key{"2014", "abc"}

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
	assert.Equal(t, &Key{"b", "c/d"}, k2)
}
