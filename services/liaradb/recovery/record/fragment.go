package record

import (
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	FragmentHeaderSize = 2
)

type Fragment struct {
	index  wrap.Int16
	buffer *buffer.Buffer
}

func NewFragment(header []byte, data []byte) Fragment {
	index, _ := wrap.NewInt16(header)
	return Fragment{
		index:  index,
		buffer: buffer.NewFromSlice(data),
	}
}

func (f Fragment) Length() int64 { return f.buffer.Length() }
func (f Fragment) Index() int16  { return f.index.Get() }

func (f Fragment) SetIndex(v int16) {
	f.index.Set(v)
}

func (f Fragment) Read(p []byte) (n int, err error) {
	return f.buffer.Read(p)
}

func (f Fragment) Write(p []byte) (n int, err error) {
	return f.buffer.Write(p)
}
