package finch

import "strings"
import "errors"

var ErrDeserializeKey error = errors.New("not enough parts to split")

const (
	TasksIndex     string = "tasks"
	SelectedIndex  string = "selected"
	AvailableIndex string = "available"
)

type Key struct {
	Index     string
	Timestamp string
	Hash      string
}

func DeserializeKey(szd []byte) (*Key, error) {
	k := new(Key)

	parts := strings.SplitN(string(szd), "/", 3)
	if len(parts) < 3 {
		return k, ErrDeserializeKey
	}
	k.Index = parts[0]
	k.Timestamp = parts[1]
	k.Hash = parts[2]

	return k, nil
}

func (k *Key) Serialize() []byte {
	return []byte(k.Index + "/" + k.Timestamp + "/" + k.Hash)
}
