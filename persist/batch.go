package persist

import (
	"github.com/vmihailenco/msgpack"

	"time"
)

type PutDeleter interface {
	Put([]byte, []byte) error
	Delete([]byte) error
}

type LoggedOperation struct {
	Operation string
	Key       []byte
	Value     []byte
	Timestamp time.Time
}

func NewLoggedOperation(operation string, key, value []byte) *LoggedOperation {
	return &LoggedOperation{
		Operation: operation,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}
}

func (lo *LoggedOperation) LogKey() []byte {
	return []byte("_log/" + lo.Timestamp.Format(time.RFC3339))
}

func (lo *LoggedOperation) Serialize() ([]byte, error) {
	return msgpack.Marshal(lo)
}

type LoggedBatch struct {
	Batch PutDeleter
}

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
