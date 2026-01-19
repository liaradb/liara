package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

type Page interface {
	Add([]byte) (page.Offset, error)
	ID() action.PageID
	TimeLineID() action.TimeLineID
	LengthRemaining() record.Length
	Items() iter.Seq2[[]byte, error]
	ItemsReverse() iter.Seq2[[]byte, error]
	Read(r io.ReadSeeker) error
	Reset(action.PageID, action.TimeLineID, record.Length)
	Write(w io.WriteSeeker) error
}
