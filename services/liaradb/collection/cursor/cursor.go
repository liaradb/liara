package cursor

import (
	"io"

	"github.com/liaradb/liaradb/encoder/buffer"
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

func (c *Cursor) Writer() io.Writer {
	writers := make([]io.Writer, len(c.buffers))
	for _, b := range c.buffers {
		writers = append(writers, buffer.NewFromSlice(b.Raw()))
	}
	return io.MultiWriter(writers...)
}
