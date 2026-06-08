package record

import (
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	FragmentHeaderSize = 4
)

type Fragment struct {
	count  wrap.Int16
	index  wrap.Int16
	buffer *buffer.Buffer
}

func NewFragment(header []byte, data []byte) Fragment {
	count, header0 := wrap.NewInt16(header)
	index, _ := wrap.NewInt16(header0)
	return Fragment{
		count:  count,
		index:  index,
		buffer: buffer.NewFromSlice(data),
	}
}

func (f Fragment) Length() int64 { return f.buffer.Length() }
func (f Fragment) Count() int16  { return f.count.Get() }
func (f Fragment) Index() int16  { return f.index.Get() }

func (f Fragment) SetCount(v int16) {
	f.count.Set(v)
}

func (f Fragment) SetIndex(v int16) {
	f.index.Set(v)
}

func (f Fragment) Read(p []byte) (n int, err error) {
	return f.buffer.Read(p)
}

func (f Fragment) Write(p []byte) (n int, err error) {
	return f.buffer.Write(p)
}
