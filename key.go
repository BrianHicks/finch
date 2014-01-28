package finch

import (
	"strings"
	"time"
)
import "errors"

var (
	// ErrDeserializeKey is returned when a key is too short (doesn't have
	// enough parts)
	ErrDeserializeKey = errors.New("not enough parts to split")
)

const (
	// PendingIndex is the string key for pending tasks
	PendingIndex string = "pending"
	// SelectedIndex is the string key for selected tasks
	SelectedIndex string = "selected"
	// TasksIndex is the prefix all keys are stored under
	TasksIndex string = "tasks"
)

// Key serializes and deserializes ordered key information in LevelDB
type Key struct {
	Index     string
	Timestamp string
	ID        string
}

// KeyForTask returns a Key for a Task
func KeyForTask(idx string, t *Task) *Key {
	return &Key{
		idx,
		t.Added.Format(time.RFC3339),
		t.ID,
	}
}

// DeserializeKey from a []byte. It tries to deal with errors gracefully but
// can't if it doesn't have enough input.
func DeserializeKey(szd []byte) (*Key, error) {
	k := new(Key)

	parts := strings.SplitN(string(szd), "/", 3)
	if len(parts) < 3 {
		return k, ErrDeserializeKey
	}
	k.Index = parts[0]
	k.Timestamp = parts[1]
	k.ID = parts[2]

	return k, nil
}

// Serialize a key to a []byte by adding "/" between it's fields.
//
// A constructed key looks like "tasks/2014-01-01T00:00:00/<sha1>"
func (k *Key) Serialize() []byte {
	return []byte(k.Index + "/" + k.Timestamp + "/" + k.ID)
}
