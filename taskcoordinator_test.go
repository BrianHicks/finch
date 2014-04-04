package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskCoordinatorAdd(t *testing.T) {
	t.Parallel()

	store, err := NewJSONStore("tcadd.json")
	assert.Nil(t, err)

	tc := TaskCoordinator{store}
	task, err := tc.Add("test")
	assert.Nil(t, err)
	assert.Equal(t, task.Desc, "test")
}
