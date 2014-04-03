package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextID(t *testing.T) {
	t.Parallel()

	j := JSONStore{CurID: 0}

	i := j.NextID()
	assert.Equal(t, i, 1)
	assert.Equal(t, j.CurID, 1)
}

func TestCommit(t *testing.T) {
	j := JSONStore{filename: "test.json", CurID: 0}
	fname := "test.json"
	defer os.Remove(fname)

	err := j.Commit()
	assert.Nil(t, err)

	bytes, err := ioutil.ReadFile(fname)
	assert.Nil(t, err)
	assert.True(t, len(bytes) > 0)
}
