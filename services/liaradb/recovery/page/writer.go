package page

import (
	"io"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/node"
	"github.com/liaradb/liaradb/recovery/record"
)

type Writer struct {
	bodySize int64
	page     *node.Node
}

func NewWriter(size int64, page *node.Node) *Writer {
	return &Writer{
		bodySize: size,
		page:     page,
	}
}

func (wr *Writer) Init(id action.PageID, tlid action.TimeLineID, rem record.Length) {
	wr.page.Reset(id, tlid, rem)
}

func (wr *Writer) Append(data []byte) bool {
	_, ok := wr.page.Append(data)
	return ok
}

func (wr *Writer) Position() int64 {
	return wr.page.ID().Position(wr.bodySize)
}

func (wr *Writer) Write(w io.WriterAt) error {
	return wr.page.Write(io.NewOffsetWriter(w, wr.Position()))
}

func (wr *Writer) Read(r io.ReadSeeker) error {
	return wr.page.Read(r)
}
