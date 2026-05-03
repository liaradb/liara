package record

import (
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/storage/link"
)

type Fragment struct {
	buffer *buffer.Buffer
	first  bool
	next   link.RecordLocator
	prev   link.RecordLocator
}

func NewFragment(data []byte) Fragment {
	return Fragment{
		buffer: buffer.NewFromSlice(data),
	}
}

func (f Fragment) Length() int64 { return f.buffer.Length() }

func (f Fragment) Read(p []byte) (n int, err error) {
	return f.buffer.Read(p)
}

func (f Fragment) Write(p []byte) (n int, err error) {
	return f.buffer.Write(p)
}
