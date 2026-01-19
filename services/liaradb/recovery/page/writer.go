package page

import (
	"io"

	"github.com/liaradb/liaradb/recovery/record"
)

type Writer struct {
	bodySize int64
	page     Page
}

func NewWriter(size int64) *Writer {
	return &Writer{
		bodySize: size,
		page:     newPage(size),
	}
}

func (wr *Writer) Init(id PageID, tlid TimeLineID, rem record.Length) {
	// TODO: Don't replace header
	wr.page.Reset(NewHeader(id, tlid, rem))
}

func (wr *Writer) Append(data []byte) error {
	_, err := wr.page.Add(data)
	return err
}

func (wr *Writer) Position() int64 {
	return wr.page.Header().ID().Position(wr.bodySize)
}

func (wr *Writer) Write(w io.WriterAt) error {
	return wr.page.Write(io.NewOffsetWriter(w, wr.Position()))
}

func (wr *Writer) Read(r io.ReadSeeker) error {
	return wr.page.Read(r)
}
