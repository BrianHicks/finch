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

func TestJSONStoreFilterTasks(t *testing.T) {
	t.Parallel()
	j, err := NewJSONStore("filter.json")
	assert.Nil(t, err)

	foo := &Task{ID: "foo"}
	bar := &Task{ID: "bar"}
	j.SaveTask(foo)
	j.SaveTask(bar)

	// accept everything
	all := func(t *Task) bool { return true }

	result, err := j.FilterTasks(all)
	assert.Nil(t, err)
	in := func(t *testing.T, task *Task, cont []*Task) {
		for _, x := range cont {
			if x.ID == task.ID {
				return
			}
		}
		t.Errorf("Task with id %s not found in %+v", task.ID, cont)
	}
	in(t, foo, result)
	in(t, bar, result)

	// only accept Task with ID "foo"
	some := func(t *Task) bool { return t.ID == "foo" }

	result, err = j.FilterTasks(some)
	assert.Nil(t, err)
	assert.Equal(t, []*Task{foo}, result)
}

func TestJSONStoreDeleteTask(t *testing.T) {
	t.Parallel()

	j, err := NewJSONStore("delete.json")
	assert.Nil(t, err)

	task := Task{ID: "foo", Desc: "test"}
	j.Tasks[task.ID] = &task

	// deleting none is an error
	err = j.DeleteTask("")
	assert.Equal(t, err, NoSuchTask)

	// deleting one should work
	err = j.DeleteTask(task.ID)
	assert.Nil(t, err)
	_, ok := j.Tasks[task.ID]
	assert.False(t, ok)
}

func TestJSONStoreSetMeta(t *testing.T) {
	t.Parallel()
	j, err := NewJSONStore("setmeta.json")
	assert.Nil(t, err)

	k, v := "foo", "bar"
	err = j.SetMeta(k, v)
	assert.Nil(t, err)

	// make sure it actually set
	v2, ok := j.Meta[k]
	assert.True(t, ok)
	assert.Equal(t, v, v2)
}

func TestJSONStoreGetMeta(t *testing.T) {
	t.Parallel()
	j, err := NewJSONStore("getmeta.json")
	assert.Nil(t, err)

	k, v := "foo", "bar"
	j.Meta[k] = v

	// get a value that exists
	v2, err := j.GetMeta(k)
	assert.Nil(t, err)
	assert.Equal(t, v, v2)

	// get a value that doesn't exist
	_, err = j.GetMeta("whatever")
	assert.Equal(t, err, NoSuchKey)
}

func TestJSONStoreImplements(t *testing.T) {
	t.Parallel()

	j := new(JSONStore)
	assert.Implements(t, (*Storage)(nil), j)
	assert.Implements(t, (*TaskStore)(nil), j)
	assert.Implements(t, (*MetaStore)(nil), j)
	assert.Implements(t, (*MetaTaskStore)(nil), j)
}
