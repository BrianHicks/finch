package finch

import "strings"

import "errors"

var (
	// ErrDeserializeKey is returned when a key is too short (doesn't have
	// enough parts)
	ErrDeserializeKey = errors.New("not enough parts to split")
)

const (
	// TasksIndex is the prefix all keys are stored under
	TasksIndex string = "tasks"
)

// Key serializes and deserializes ordered key information in LevelDB
type Key struct {
	Timestamp string
	ID        string
}

// DeserializeKey from a []byte. It tries to deal with errors gracefully but
// can't if it doesn't have enough input.
func DeserializeKey(szd []byte) (*Key, error) {
	k := new(Key)

	parts := strings.SplitN(string(szd), "/", 3)
	if len(parts) < 3 {
		return k, ErrDeserializeKey
	}
	k.Timestamp = parts[1]
	k.ID = parts[2]

	return k, nil
}

// Serialize a key to a []byte by adding "/" between it's fields. It will be
// serialized for the prefix given.
//
// A constructed key looks like "tasks/2014-01-01T00:00:00/<sha1>"
func (k *Key) Serialize(prefix string) []byte {
	return []byte(prefix + "/" + k.Timestamp + "/" + k.ID)
}
