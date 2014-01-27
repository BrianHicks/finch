package finch

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestSerializationRoundtrip(t *testing.T) {
	origTask := new(Task)
	origTask.Description = "test!"
	origTask.Pending = true
	origTask.Added = time.Now()

	serialized, err := origTask.Serialize()
	assert.Nil(t, err)

	newTask, err := TaskFromSerialized(serialized)
	assert.Nil(t, err)

	assert.Equal(t, newTask, origTask)
}

func TestNewTask(t *testing.T) {
	now := time.Now()
	task := NewTask("test!", now)

	assert.Equal(t, task.Description, "test!")
	assert.Equal(t, task.Added, now)
	assert.False(t, task.Selected)
	assert.True(t, task.Pending)
}
