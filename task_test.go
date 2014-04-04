package main

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskString(t *testing.T) {
	t.Parallel()

	// selected
	assert.Equal(
		t,
		(&Task{ID: "foo", Desc: "test", Selected: true}).String(),
		"foo: test (*)",
	)

	// unselected
	assert.Equal(
		t,
		(&Task{ID: "foo", Desc: "test", Selected: false}).String(),
		"foo: test",
	)

	// done
	assert.Equal(
		t,
		(&Task{ID: "foo", Desc: "test", Done: true}).String(),
		"foo: test (done)",
	)
}

func TestByActive(t *testing.T) {
	t.Parallel()

	now := time.Now()
	past := &Task{ID: "past", Active: now.Add(-1 * time.Second)}
	present := &Task{ID: "present", Active: now}
	future := &Task{ID: "future", Active: now.Add(time.Second)}

	tasks := []*Task{present, future, past}

	sort.Sort(ByActive(tasks))

	assert.Equal(
		t,
		tasks,
		[]*Task{past, present, future},
	)
}

func TestTaskMarkDone(t *testing.T) {
	t.Parallel()

	// task with no Repeat
	task := Task{Selected: true}
	assert.False(t, task.Done)
	assert.True(t, task.Selected)

	task.MarkDone()
	assert.True(t, task.Done)
	assert.False(t, task.Selected)

	// task with a Repeat
	now := time.Now()
	task = Task{Selected: true, Repeat: time.Second}
	task.MarkDone()

	assert.False(t, task.Done)
	assert.True(t, task.Active.After(now))
}
