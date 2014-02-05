package persist

type KV struct {
	Key, Value []byte
}

type KeyValuer interface {
	Key() []byte
	Value() []byte
}

func NewKVFromIter(iter KeyValuer) *KV {
	kv := new(KV)

	key := iter.Key()
	kv.Key = make([]byte, len(key))
	copy(kv.Key, key)

	value := iter.Value()
	kv.Value = make([]byte, len(value))
	copy(kv.Value, value)

	return kv
}
