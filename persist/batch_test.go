package persist

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vmihailenco/msgpack"

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

type MockedBatch struct {
	mock.Mock
}

func (m *MockedBatch) Put(key, value []byte) error {
	// only compare the first 30 bytes of value since timestamps differ and
	// that's always at the end.
	if len(value) > 30 {
		value = value[:30]
	}
	args := m.Mock.Called(key, value)
	return args.Error(0)
}

func (m *MockedBatch) Delete(key []byte) error {
	args := m.Mock.Called(key)
	return args.Error(0)
}

func TestLoggedBatchPut(t *testing.T) {
	t.Parallel()

	key := []byte{1}
	value := []byte{2}

	logged := NewLoggedOperation("PUT", key, value)
	szd, err := logged.Serialize()
	assert.Nil(t, err)

	mockBatch := new(MockedBatch)
	mockBatch.On("Put", key, value).Return(nil)
	mockBatch.On("Put", logged.LogKey(), szd[:30]).Return(nil)

	batch := LoggedBatch{mockBatch}
	batch.Put(key, value)

	mockBatch.Mock.AssertExpectations(t)
}

func TestLoggedBatchDelete(t *testing.T) {
	t.Parallel()

	key := []byte{1}

	logged := NewLoggedOperation("DELETE", key, nil)
	szd, err := logged.Serialize()
	assert.Nil(t, err)

	mockBatch := new(MockedBatch)
	mockBatch.On("Delete", key).Return(nil)
	mockBatch.On("Put", logged.LogKey(), szd[:30]).Return(nil)

	batch := LoggedBatch{mockBatch}
	batch.Delete(key)

	mockBatch.Mock.AssertExpectations(t)
}
