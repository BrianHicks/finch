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

	task := tc.Add("test")
	tc.Save(task)

	assert.Equal(t, task.Desc, "test")
}

func TestTaskCoordinatorDelay(t *testing.T) {
	tc := setup(t, "tcdelay.json")

	task := tc.Add("test")
	tc.Save(task)

	now := task.Active.Add(time.Second)

	// no ID
	err := tc.Delay("", now)
	assert.Equal(t, err, NoSuchTask)

	// with ID
	assert.NotEqual(t, task.Active, now)

	err = tc.Delay(task.ID, now)
	assert.Nil(t, err)
	assert.Equal(t, task.Active, now)
}

func TestTaskCoordinatorSelect(t *testing.T) {
	tc := setup(t, "tcselect.json")

	task := tc.Add("test")
	tc.Save(task)

	// no ID
	err := tc.Select("")
	assert.Equal(t, err, NoSuchTask)

	// with ID
	err = tc.Select(task.ID)
	assert.Nil(t, err)
	assert.True(t, task.Selected)
}

func TestTaskCoordinatorMarkDone(t *testing.T) {
	tc := setup(t, "tcmarkdone.json")

	task := tc.Add("test")
	tc.Save(task)

	// no ID
	err := tc.MarkDone("")
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
	task := tc.Add("test")
	tc.Save(task)

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
	task := tc.Add("test")
	task.Selected = true
	tc.Save(task)

	task2, err := tc.NextSelected()
	assert.Nil(t, err)
	assert.Equal(t, task, task2)
}

func TestTaskCoordinatorAvailable(t *testing.T) {
	tc := setup(t, "tcavailable.json")

	// no available tasks should return an error
	tasks, err := tc.Available()
	assert.Equal(t, err, NoSuchTask)
	assert.Equal(t, len(tasks), 0)

	// some tasks
	done := tc.Add("done")
	done.Done = true
	tc.Save(done)

	future := tc.Add("future")
	future.Active = time.Now().Add(time.Second * 30)
	tc.Save(future)

	pending := tc.Add("pending")
	tc.Save(pending)

	tasks, err = tc.Available()
	assert.Nil(t, err)
	assert.Equal(t, []*Task{pending}, tasks)
}
