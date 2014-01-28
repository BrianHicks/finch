package finch

import (
	"strings"
	"time"
)
import "errors"

var ErrDeserializeKey error = errors.New("not enough parts to split")

const (
	PendingIndex  string = "pending"
	SelectedIndex string = "selected"
	TasksIndex    string = "tasks"
)

type Key struct {
	Index     string
	Timestamp string
	Id        string
}

func KeyForTask(idx string, t *Task) *Key {
	return &Key{
		idx,
		t.Added.Format(time.RFC3339),
		t.Id,
	}
}

func DeserializeKey(szd []byte) (*Key, error) {
	k := new(Key)

	parts := strings.SplitN(string(szd), "/", 3)
	if len(parts) < 3 {
		return k, ErrDeserializeKey
	}
	k.Index = parts[0]
	k.Timestamp = parts[1]
	k.Id = parts[2]

	return k, nil
}

func (k *Key) Serialize() []byte {
	return []byte(k.Index + "/" + k.Timestamp + "/" + k.Id)
}
