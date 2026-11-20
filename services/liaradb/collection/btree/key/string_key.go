package key

type stringKey string

func (k stringKey) Size() int      { return len(k) }
func (k stringKey) String() string { return string(k) }
func (k stringKey) Value() []byte  { return []byte(k) }

func (k stringKey) Equal(b stringKey) bool {
	return k == b
}

func (k stringKey) Greater(b stringKey) bool {
	return k > b
}

func (k stringKey) GreaterEqual(b stringKey) bool {
	return k >= b
}

func (k stringKey) Less(b stringKey) bool {
	return k < b
}

func (k stringKey) LessEqual(b stringKey) bool {
	return k <= b
}
