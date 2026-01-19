package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
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
