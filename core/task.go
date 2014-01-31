package core

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/vmihailenco/msgpack"

	"io"
	"time"
)

const (
	// TagPending is used to annotate tasks which haven't been completed yet
	TagPending = "pending"
	// TagSelected shows selection status - IE as part of a FV chain
	TagSelected = "selected"
)

// Task is the basic structure for tasks in the database
type Task struct {
	ID          string
	Description string
	Timestamp   time.Time
	Attrs       map[string]bool
}

// NewTask returns an instantiated Task. In particular, it hashes the
// description into "ID"
func NewTask(description string, added time.Time) *Task {
	t := new(Task)
	t.Description = description
	t.Attrs = map[string]bool{
		TagPending:  true,
		TagSelected: false,
	}
	t.Timestamp = added

	hash := sha1.New()
	io.WriteString(hash, description)
	t.ID = base64.StdEncoding.EncodeToString(hash.Sum(nil))

	return t
}

// DeserializeTask from a []byte, assumes msgpack encoding
func DeserializeTask(szd []byte) (*Task, error) {
	t := new(Task)
	err := msgpack.Unmarshal(szd, t)

	return t, err
}

// Serialize a task to a msgpack-encoded []byte
func (t *Task) Serialize() ([]byte, error) {
	b, err := msgpack.Marshal(t)
	return b, err
}

// Key returns a valid core.Key for this Task
func (t *Task) Key() *Key {
	return &Key{
		Timestamp: t.Timestamp.Format(time.RFC3339),
		ID:        t.ID,
	}
}
