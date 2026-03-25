package key

import (
	"fmt"

	"github.com/liaradb/liaradb/encoder/scan"
)

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

func (k Key) Size() int      { return k.A.Size() + k.B.Size() }
func (k Key) String() string { return fmt.Sprintf("%v%v", k.A, k.B) }

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

func (k Key) Write(data []byte) bool {
	data0, ok := scan.SetInt64(data, k.B.Value())
	if !ok {
		return false
	}

	copy(data0, k.A.Value())
	return true
}

func (k *Key) Read(data []byte) bool {
	b, data0, ok := scan.Int64(data)
	if !ok {
		return false
	}

	k.B = intKey(b)
	k.A = stringKey(data0)
	return true
}
