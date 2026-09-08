package span

import (
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/logpage"
)

const (
	FragmentHeaderSize = base.Uint16Size +
		base.Uint16Size +
		page.CrcSize
)

type Fragment struct {
	l      Log
	p      BufferPage
	count  wrap.Int16
	index  wrap.Int16
	crc    wrap.Int32
	buffer *buffer.Buffer
}

type BufferPage interface {
	SetLogSequenceNumber(logpage.LogSequenceNumber)
}

func newFragment(
	l Log,
	p BufferPage,
	header []byte,
	data []byte,
) *Fragment {
	count, header0 := wrap.NewInt16(header)
	index, header1 := wrap.NewInt16(header0)
	crc, _ := wrap.NewInt32(header1)
	return &Fragment{
		l:      l,
		p:      p,
		count:  count,
		index:  index,
		crc:    crc,
		buffer: buffer.NewFromSlice(data),
	}
}

// TODO: Fix this cast
func (f Fragment) length() int  { return int(f.buffer.Length()) }
func (f Fragment) Count() int16 { return f.count.Get() }
func (f Fragment) Index() int16 { return f.index.Get() }

func (f Fragment) valid() bool {
	return page.RestoreCRC(f.crc.Get()).
		Compare(f.buffer.Bytes())
}

func (f Fragment) commit() {
	crc := page.NewCRC(f.buffer.Bytes())
	f.crc.Set(int32(crc.Value()))
}

func (f Fragment) setCount(v int16) {
	f.count.Set(v)
}

func (f Fragment) setIndex(v int16) {
	f.index.Set(v)
}

func (f Fragment) Read(p []byte) (int, error) {
	return f.buffer.Read(p)
}

func (f Fragment) Write(p []byte) (int, error) {
	return f.buffer.Write(p)
}
