package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
)

type Page[I ItemSerializer] interface {
	Add(i I) (page.Offset, error)
	Header() *Header
	Items() iter.Seq2[I, error]
	ItemsReverse() iter.Seq2[I, error]
	Read(r io.ReadSeeker) error
	Reset(h *Header)
	Size() int
	Write(w io.WriteSeeker) error
}

type ItemSerializer interface {
	Read(io.Reader, page.CRC) error
	Size() int
	Write(io.Writer) (page.CRC, error)
}
