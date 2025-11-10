package page

import (
	"io"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Writer struct {
	bodySize int64
	page     *page.Page[*Header, *page.Item]
}

func NewWriter(size int64) *Writer {
	return &Writer{
		bodySize: size,
		page: page.NewWithHeader(
			page.Offset(size),
			&Header{},
			page.NewItemByLength),
	}
}

func (wr *Writer) Init(id PageID, tlid TimeLineID, rem record.Length) {
	h := NewHeader(id, tlid, rem)
	// TODO: Don't replace page
	wr.page = page.NewWithHeader(
		page.Offset(wr.bodySize),
		&h,
		page.NewItemByLength)
}

func (wr *Writer) Append(data []byte) error {
	_, err := wr.page.Add(page.NewItem(data))
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
