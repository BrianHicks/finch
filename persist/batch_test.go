package persist

import (
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"
	"strings"

	"testing"
	"time"
)

func TestLogKey(t *testing.T) {
	t.Parallel()

	now := time.Now()
	op := LoggedOperation{
		Operation: "TEST",
		Key:       nil,
		Value:     nil,
		Timestamp: now,
	}

	assert.Equal(
		t,
		[]byte("_log/"+now.Format(time.RFC3339)),
		op.LogKey(),
	)
}

func TestLoggedOperationSerialize(t *testing.T) {
	t.Parallel()

	op := LoggedOperation{"TEST", nil, nil, time.Now()}
	serialized, err := op.Serialize()
	assert.Nil(t, err)

	packed, err := msgpack.Marshal(op)
	assert.Nil(t, err)

	assert.Equal(t, packed, serialized)
}

func TestLoggedBatchPut(t *testing.T) {
	t.Parallel()
	db, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte{1}
	value := []byte{2}

	batch := LoggedBatch{&leveldb.Batch{}}
	batch.Put(key, value)

	err = db.Write(batch.Batch, db.WO)
	assert.Nil(t, err)

	// check document
	doc, err := db.Get(key, db.RO)
	assert.Nil(t, err)
	assert.Equal(t, value, doc)

	// check log
	iter := db.NewIterator(db.RO)
	ok := iter.Seek([]byte("_log/"))
	assert.True(t, ok)
	assert.True(t, strings.Contains(string(iter.Key()), "_log/"))
	assert.True(t, strings.Contains(string(iter.Value()), "PUT"))
}

func TestLoggedBatchDelete(t *testing.T) {
	t.Parallel()
	db, err := NewInMemory()
	assert.Nil(t, err)

	key := []byte{1}

	batch := LoggedBatch{&leveldb.Batch{}}
	batch.Delete(key)

	err = db.Write(batch.Batch, db.WO)
	assert.Nil(t, err)

	// check document
	_, err = db.Get(key, db.RO)
	assert.NotNil(t, err)

	// check log
	iter := db.NewIterator(db.RO)
	ok := iter.Seek([]byte("_log/"))
	assert.True(t, ok)
	assert.True(t, strings.Contains(string(iter.Key()), "_log/"))
	assert.True(t, strings.Contains(string(iter.Value()), "DELETE"))
}
