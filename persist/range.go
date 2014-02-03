package persist

// Range implements a range over some LevelDB keys. Since this always has to
// know what DB to range over, use DB.Key() and similar to construct.
type Range struct {
	Start []byte
	End   []byte
	db    *DB
}

func (r *Range) First() []byte {
	return []byte{}
}

func (r *Range) Last() []byte {
	return []byte{}
}

func (r *Range) All() [][]byte {
	return [][]byte{}
}
