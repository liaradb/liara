package btree

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Key string

// TODO: Test this
func (i Key) String() string {
	return string(i)
}

func (k Key) Length() int16 {
	return int16(len(k))
}

func (k Key) Equal(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k == b
}

func (k Key) Greater(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k > b
}

func (k Key) GreaterEqual(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k >= b
}

func (k Key) Less(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k < b
}

func (k Key) LessEqual(o any) bool {
	b, ok := o.(Key)
	if !ok {
		return false
	}

	return k <= b
}

// TODO: Test this
func (i Key) Size() int               { return raw.StringSize(i) }
func (i Key) Write(w io.Writer) error { return raw.WriteString(w, i) }
func (i *Key) Read(r io.Reader) error { return raw.ReadString(r, i) }
