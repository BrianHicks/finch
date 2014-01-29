package finch

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestSerializationRoundtrip(t *testing.T) {
	origTask := new(Task)
	origTask.Description = "test!"
	origTask.Added = time.Now()
	origTask.Attrs = map[string]bool{}

	serialized, err := origTask.Serialize()
	assert.Nil(t, err)

	newTask, err := DeserializeTask(serialized)
	assert.Nil(t, err)

	assert.Equal(t, newTask, origTask)
}

func TestNewTask(t *testing.T) {
	now := time.Now()
	task := NewTask("test!", now)

	assert.Equal(t, task.Description, "test!")
	assert.Equal(t, task.Added, now)
	assert.False(t, task.Attrs[TagSelected])
	assert.True(t, task.Attrs[TagPending])
}

func TestKey(t *testing.T) {
	task := NewTask("test", time.Now())

	key := task.Key()

	assert.Equal(t, task.Added.Format(time.RFC3339), key.Timestamp)
	assert.Equal(t, task.ID, key.ID)
}
