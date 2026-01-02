package value

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type key string

// TODO: Test this
func (i key) String() string {
	return string(i)
}

func (k key) Value() []byte {
	return []byte(k)
}

func (k key) Length() int16 {
	return int16(len(k))
}

func (k key) Equal(o any) bool {
	b, ok := o.(key)
	if !ok {
		return false
	}

	return k == b
}

func (k key) Greater(o any) bool {
	b, ok := o.(key)
	if !ok {
		return false
	}

	return k > b
}

func (k key) GreaterEqual(o any) bool {
	b, ok := o.(key)
	if !ok {
		return false
	}

	return k >= b
}

func (k key) Less(o any) bool {
	b, ok := o.(key)
	if !ok {
		return false
	}

	return k < b
}

func (k key) LessEqual(o any) bool {
	b, ok := o.(key)
	if !ok {
		return false
	}

	return k <= b
}

// TODO: Test this
func (i key) Size() int               { return len(i) }
func (i key) Write(w io.Writer) error { return raw.WriteString(w, i) }
func (i *key) Read(r io.Reader) error { return raw.ReadString(r, i) }
