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

func (k Key) Equal(b Key) bool {
	return k.A.Equal(b.A) && k.B.Equal(b.B)
}

func (k Key) Greater(b Key) bool {
	return k.A.Greater(b.A) || (k.A.Equal(b.A) && k.B.Greater(b.B))
}

func (k Key) GreaterEqual(b Key) bool {
	return k.A.Greater(b.A) || (k.A.Equal(b.A) && k.B.GreaterEqual(b.B))
}

func (k Key) Less(b Key) bool {
	return k.A.Less(b.A) || (k.A.Equal(b.A) && k.B.Less(b.B))
}

func (k Key) LessEqual(b Key) bool {
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
