package main

import (
	"github.com/BrianHicks/finch/core"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
	"time"
)

func TestDelayer(t *testing.T) {
	os.Setenv("CORE_STORAGE", "mem")
	tdb, err := getTaskDB()
	assert.Nil(t, err)

	task := core.NewTask("test", time.Now().Add(-1*(time.Second*30)))
	task.Attrs[core.TagSelected] = true

	err = tdb.PutTasks(task)
	assert.Nil(t, err)

	updated, err := Delayer(tdb, []string{})
	assert.Nil(t, err)

	assert.NotEqual(t, task, updated)

	// get old and make sure it doesn't exist
	_, err = tdb.GetTask(task.Key())
	assert.Equal(t, core.ErrNoTask, err)

	// get old and make sure it doesn't error
	updated_check, err := tdb.GetTask(updated.Key())
	assert.Nil(t, err)
	assert.Equal(t, updated, updated_check)
}
