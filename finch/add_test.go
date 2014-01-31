package main

import (
	commander "code.google.com/p/go-commander"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	os.Setenv("FINCH_STORAGE", "mem")

	Add.Run(&commander.Command{}, []string{"this is a test"})
	// TODO: still trying to figure out how to test these.
}

func TestAdder(t *testing.T) {
	os.Setenv("FINCH_STORAGE", "mem")

	tdb, err := getTaskDB()
	assert.Nil(t, err)

	task, err := Adder(tdb, []string{"test"})
	assert.Nil(t, err)

	confirm, err := tdb.GetTask(task.Key())
	assert.Equal(t, task, confirm)
}
