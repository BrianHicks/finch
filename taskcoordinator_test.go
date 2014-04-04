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

func TestTaskCoordinatorSelect(t *testing.T) {
	tc := setup(t, "tcselect.json")

	task, err := tc.Add("test")
	assert.Nil(t, err)

	// no ID
	err = tc.Select("")
	assert.Equal(t, err, NoSuchTask)

	// with ID
	err = tc.Select(task.ID)
	assert.Nil(t, err)
	assert.True(t, task.Selected)
}

func TestTaskCoordinatorMarkDone(t *testing.T) {
	tc := setup(t, "tcmarkdone.json")

	task, err := tc.Add("test")
	assert.Nil(t, err)

	// no ID
	err = tc.MarkDone("")
	assert.Equal(t, err, NoSuchTask)

	// with ID
	err = tc.MarkDone(task.ID)
	assert.Nil(t, err)
	assert.True(t, task.Done)
}

func TestTaskCoordinatorSelected(t *testing.T) {
	tc := setup(t, "tcselected.json")

	// no tasks gets an error
	tasks, err := tc.Selected()
	assert.Equal(t, err, NoSuchTask)
	assert.Equal(t, len(tasks), 0)

	// some selected tasks returns those
	task, err := tc.Add("test")
	assert.Nil(t, err)

	task.Selected = true

	tasks, err = tc.Selected()
	assert.Nil(t, err)
	assert.Equal(t, tasks, []*Task{task})
}

func TestTaskCoordinatorNextSelected(t *testing.T) {
	tc := setup(t, "tcnextselected.json")

	// no tasks gets an error
	notask, err := tc.NextSelected()
	assert.Equal(t, err, NoSuchTask)
	assert.Nil(t, notask)

	// some selected tasks returns that one
	task, err := tc.Add("test")
	assert.Nil(t, err)
	task.Selected = true

	task2, err := tc.NextSelected()
	assert.Nil(t, err)
	assert.Equal(t, task, task2)
}
