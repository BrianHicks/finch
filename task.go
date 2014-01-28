package finch

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/vmihailenco/msgpack"

	"io"
	"time"
)

type Task struct {
	Id          string
	Description string
	Added       time.Time
	Selected    bool
	Pending     bool
}

func NewTask(description string, added time.Time) *Task {
	t := new(Task)
	t.Description = description
	t.Added = added
	t.Selected = false
	t.Pending = true

	hash := sha1.New()
	io.WriteString(hash, description)
	t.Id = base64.StdEncoding.EncodeToString(hash.Sum(nil))

	return t
}

func DeserializeTask(szd []byte) (*Task, error) {
	t := new(Task)
	err := msgpack.Unmarshal(szd, t)

	return t, err
}

func (t *Task) Serialize() ([]byte, error) {
	b, err := msgpack.Marshal(t)
	return b, err
}

func (t *Task) Key(idx string) *Key {
	return KeyForTask(idx, t)
}
