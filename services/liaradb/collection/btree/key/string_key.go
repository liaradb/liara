package key

type stringKey string

func (k stringKey) Size() int      { return len(k) }
func (k stringKey) String() string { return string(k) }
func (k stringKey) Value() []byte  { return []byte(k) }

func (k stringKey) Equal(o any) bool {
	b, ok := o.(stringKey)
	if !ok {
		return false
	}

	return k == b
}

func (k stringKey) Greater(o any) bool {
	b, ok := o.(stringKey)
	if !ok {
		return false
	}

	return k > b
}

func (k stringKey) GreaterEqual(o any) bool {
	b, ok := o.(stringKey)
	if !ok {
		return false
	}

	return k >= b
}

func (k stringKey) Less(o any) bool {
	b, ok := o.(stringKey)
	if !ok {
		return false
	}

	return k < b
}

func (k stringKey) LessEqual(o any) bool {
	b, ok := o.(stringKey)
	if !ok {
		return false
	}

	return k <= b
}
