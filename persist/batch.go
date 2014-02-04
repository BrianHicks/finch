package persist

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"

	"time"
)

// LoggedOperation encodes an operation that took place on the database, for
// later syncing.
type LoggedOperation struct {
	Operation string
	Key       []byte
	Value     []byte
	Timestamp time.Time
}

// NewLoggedOperation returns a timestamped LoggedOperation object from the
// arguments specifiec.
func NewLoggedOperation(operation string, key, value []byte) *LoggedOperation {
	return &LoggedOperation{
		Operation: operation,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}
}

// LogKey returns the key that this LoggedOperation will be stored under in the
// LevelDB store.
func (lo *LoggedOperation) LogKey() []byte {
	return []byte("_log/" + lo.Timestamp.Format(time.RFC3339))
}

// Serialize marshals this LoggedOperation as msgpack and returns the byte
// slice.
func (lo *LoggedOperation) Serialize() ([]byte, error) {
	return msgpack.Marshal(lo)
}

// LoggedBatch is a wrapper around leveldb.Batch that logs Put and Delete
// operations.
type LoggedBatch struct {
	Batch *leveldb.Batch
}

// Put adds the key and value to the batch operation, and logs them
func (lb *LoggedBatch) Put(key, value []byte) error {
	op := NewLoggedOperation("PUT", key, value)
	szd, err := op.Serialize()
	if err != nil {
		return err
	}

	lb.Batch.Put(key, value)
	lb.Batch.Put(op.LogKey(), szd)

	return nil
}

// Delete deletes a key from the operations and logs removal.
func (lb *LoggedBatch) Delete(key []byte) error {
	op := NewLoggedOperation("DELETE", key, nil)
	szd, err := op.Serialize()
	if err != nil {
		return err
	}

	lb.Batch.Delete(key)
	lb.Batch.Put(op.LogKey(), szd)

	return nil
}
