package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T, name string) TaskCoordinator {
	t.Parallel()

	store, err := NewJSONStore(name)
	assert.Nil(t, err)

	return TaskCoordinator{store}
}

func TestTaskCoordinatorAdd(t *testing.T) {
	tc := setup(t, "tcadd.json")

	task, err := tc.Add("test")
	assert.Nil(t, err)
	assert.Equal(t, task.Desc, "test")
}

func TestTaskCoordinatorDelay(t *testing.T) {
	tc := setup(t, "tcdelay.json")

	task, err := tc.Add("test")
	assert.Nil(t, err)

	now := task.Active.Add(time.Second)

	// no ID
	err = tc.Delay("", now)
	assert.Equal(t, err, NoSuchTask)

	// with ID
	assert.NotEqual(t, task.Active, now)

	err = tc.Delay(task.ID, now)
	assert.Nil(t, err)
	assert.Equal(t, task.Active, now)
}
