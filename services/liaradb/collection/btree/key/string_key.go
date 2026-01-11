package key

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type stringKey string

// TODO: Test this
func (k stringKey) String() string {
	return string(k)
}

func (k stringKey) Value() []byte {
	return []byte(k)
}

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

// TODO: Test this
func (k stringKey) Size() int               { return len(k) }
func (k stringKey) Write(w io.Writer) error { return raw.WriteString(w, k) }
func (k *stringKey) Read(r io.Reader) error { return raw.ReadString(r, k) }
