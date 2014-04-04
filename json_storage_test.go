package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONStoreNextID(t *testing.T) {
	t.Parallel()

	j := JSONStore{CurID: 0}

	i := j.NextID()
	assert.Equal(t, i, uint(1))
	assert.Equal(t, j.CurID, uint(1))
}

func TestJSONStoreCommit(t *testing.T) {
	fname := "test.json"
	defer os.Remove(fname)

	j, err := NewJSONStore(fname)
	assert.Nil(t, err)

	err = j.Commit()
	assert.Nil(t, err)

	bytes, err := ioutil.ReadFile(fname)
	assert.Nil(t, err)
	assert.True(t, len(bytes) > 0)
}

func TestJSONStoreSaveTask(t *testing.T) {
	t.Parallel()

	j, err := NewJSONStore("savetask.json")
	assert.Nil(t, err)
	task := Task{}

	err = j.SaveTask(&task)
	assert.Nil(t, err)
	assert.NotEqual(t, task.ID, "")

	task2 := j.Tasks[task.ID]
	assert.Equal(t, task.ID, task2.ID)
}

func TestJSONStoreGetTask(t *testing.T) {
	t.Parallel()

	j, err := NewJSONStore("gettask.json")
	assert.Nil(t, err)

	task := &Task{ID: "foo"}
	j.Tasks[task.ID] = task

	// a task that exists
	task2, err := j.GetTask(task.ID)
	assert.Nil(t, err)
	assert.Equal(t, task.ID, task2.ID)

	// a task that doesn't exist
	task3, err := j.GetTask("bar")
	assert.Equal(t, err, NoSuchTask)
	assert.Nil(t, task3)

}

func TestJSONStoreImplements(t *testing.T) {
	t.Parallel()

	var _ Storage = new(JSONStore)
	// var _ TaskStore = new(JSONStore)
}
