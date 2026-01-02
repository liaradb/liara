package value

type Key struct {
	A key
	B key
}

func NewKey(v []byte) Key {
	return Key{A: key(v)}
}

// TODO: Test this
func (k Key) String() string {
	return k.A.String() + k.B.String()
}

func (k Key) Length() int16 {
	return k.A.Length() + k.B.Length()
}

func (k Key) Size() int {
	return k.A.Size() + k.B.Size()
}

func (k Key) Value() []byte {
	return k.A.Value()
}

func (k Key) Equal(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.Equal(b.A) && k.B.Equal(b.B)
}

func (k Key) Greater(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.Greater(b.A) // && k.B.Greater(b.B)
}

func (k Key) GreaterEqual(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.GreaterEqual(b.A) // && k.B.GreaterEqual(b.B)
}

func (k Key) Less(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.Less(b.A) // && k.B.Less(b.B)
}

func (k Key) LessEqual(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.LessEqual(b.A) // && k.B.LessEqual(b.B)
}
