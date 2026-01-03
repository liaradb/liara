package value

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type intKey int64

// TODO: Test this
func (k intKey) String() string {
	return string(rune(k))
}

func (k intKey) Value() int64 {
	return int64(k)
}

func (k intKey) Length() int16 {
	return 0 // TODO: Change to 8
}

func (k intKey) Equal(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k == b
}

func (k intKey) Greater(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k > b
}

func (k intKey) GreaterEqual(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k >= b
}

func (k intKey) Less(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k < b
}

func (k intKey) LessEqual(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k <= b
}

// TODO: Test this
func (k intKey) Size() int               { return 0 } // TODO: Change to 1
func (k intKey) Write(w io.Writer) error { return raw.WriteInt64(w, k) }
func (k *intKey) Read(r io.Reader) error { return raw.ReadInt64(r, k) }
