package key

import "github.com/liaradb/liaradb/encoder/scan"

type Key struct {
	A stringKey
	B intKey
}

func NewKey(a []byte) Key {
	return Key{A: stringKey(a)}
}

func NewKey2(a []byte, b int64) Key {
	return Key{A: stringKey(a), B: intKey(b)}
}

// TODO: Test this
func (k Key) String() string {
	return k.A.String() + k.B.String()
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

	return k.A.Greater(b.A) || (k.A.Equal(b.A) && k.B.Greater(b.B))
}

func (k Key) GreaterEqual(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.Greater(b.A) || (k.A.Equal(b.A) && k.B.GreaterEqual(b.B))
}

func (k Key) Less(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.Less(b.A) || (k.A.Equal(b.A) && k.B.Less(b.B))
}

func (k Key) LessEqual(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k.A.Less(b.A) || (k.A.Equal(b.A) && k.B.LessEqual(b.B))
}

func (k Key) Write(data []byte) {
	data0 := scan.SetInt64(data, k.B.Value())
	copy(data0, k.A.Value())
}

func (k *Key) Read(data []byte) {
	b, data0 := scan.Int64(data)
	k.B = intKey(b)
	k.A = stringKey(data0)
}
