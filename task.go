package finch

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/vmihailenco/msgpack"

	"io"
	"time"
)

// Task is the basic structure for tasks in the database
type Task struct {
	ID          string
	Description string
	Added       time.Time
	Selected    bool
	Pending     bool
}

// NewTask returns an instantiated Task. In particular, it hashes the
// description into "ID"
func NewTask(description string, added time.Time) *Task {
	t := new(Task)
	t.Description = description
	t.Added = added
	t.Selected = false
	t.Pending = true

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

// Key returns a valid finch.Key for this Task
func (t *Task) Key(idx string) *Key {
	return KeyForTask(idx, t)
}
