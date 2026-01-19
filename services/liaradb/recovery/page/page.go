package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/mempage"
)

type Page interface {
	Add([]byte) (page.Offset, error)
	Header() *Header
	Items() iter.Seq2[[]byte, error]
	ItemsReverse() iter.Seq2[[]byte, error]
	Read(r io.ReadSeeker) error
	Reset(*Header)
	Write(w io.WriteSeeker) error
}

func newPage(pageSize int64) Page {
	return mempage.NewWithHeader(pageSize, &Header{})
}
