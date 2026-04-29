package cursor

import (
	"io"

	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/multi"
	"github.com/liaradb/liaradb/storage"
)

type Cursor struct {
	buffers []*storage.Buffer
}

func New(buffers ...*storage.Buffer) *Cursor {
	return &Cursor{
		buffers: buffers,
	}
}

func (c *Cursor) Release() {
	for _, b := range c.buffers {
		b.Release()
	}
}

func (c *Cursor) Reader() io.Reader {
	readers := make([]io.Reader, 0, len(c.buffers))
	for _, b := range c.buffers {
		readers = append(readers, buffer.NewFromSlice(b.Raw()))
	}
	return multi.NewReader(readers...)
}

func (c *Cursor) Writer() io.Writer {
	writers := make([]io.Writer, 0, len(c.buffers))
	for _, b := range c.buffers {
		writers = append(writers, buffer.NewFromSlice(b.Raw()))
	}
	return multi.NewWriter(writers...)
}
