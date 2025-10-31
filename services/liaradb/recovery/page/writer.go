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

func (wr *Writer) ID() PageID                     { return wr.page.Header().ID() }
func (wr *Writer) TimeLineID() TimeLineID         { return wr.page.Header().TimeLineID() }
func (wr *Writer) LengthRemaining() record.Length { return wr.page.Header().LengthRemaining() }

func (wr *Writer) Init(id PageID, tlid TimeLineID, rem record.Length) {
	h := NewHeader(id, tlid, rem)
	// TODO: Don't replace page
	wr.page = page.NewWithHeader(
		page.Offset(wr.bodySize),
		&h,
		page.NewItemByLength)
}

func (wr *Writer) Append(data []byte) error {
	return wr.page.Add(page.NewItem(data))
}

func (wr *Writer) Flush(w io.WriteSeeker) error {
	return wr.Write(w)
}

func (wr *Writer) Position() int64 {
	return wr.page.Header().ID().Size(wr.bodySize)
}

func (wr *Writer) Write(w io.WriteSeeker) error {
	return wr.page.Write(w)
}

func (wr *Writer) SeekTail(r io.ReadSeeker) error {
	// TODO: Should we handle EOF?
	if err := wr.page.Read(r); err != nil {
		if err != io.EOF {
			return err
		}
	}

	return nil
}
