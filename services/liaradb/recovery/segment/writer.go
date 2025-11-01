package segment

import (
	"bytes"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	BlockSize   = 1024
	SegmentSize = 1024
)

type Writer struct {
	pageSize    int64
	segmentSize page.PageID
	pageID      page.PageID
	timeLineID  page.TimeLineID
	writer      io.WriterAt
	recordBuf   *bytes.Buffer
	pageWriter  *page.Writer
}

type readWriterAt interface {
	io.WriterAt
	io.ReaderAt
}

func NewWriter(
	pageSize int64,
	segmentSize page.PageID,
) *Writer {
	return &Writer{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		recordBuf:   bytes.NewBuffer(nil),
		pageWriter:  page.NewWriter(pageSize),
	}
}

func (wr *Writer) PageID() page.PageID { return wr.pageID }

func (wr *Writer) Append(rc *record.Record) error {
	data, err := wr.recordToBytes(rc)
	if err != nil {
		return err
	}

	return wr.append(data)
}

func (wr *Writer) recordToBytes(rc *record.Record) ([]byte, error) {
	wr.recordBuf.Reset()
	if err := rc.Write(wr.recordBuf); err != nil {
		return nil, err
	}

	// TODO: Don't clone
	return bytes.Clone(wr.recordBuf.Bytes()), nil
}

func (wr *Writer) append(data []byte) error {
	if err := wr.appendOrNext(data); err != nil {
		return err
	}

	return nil
}

func (wr *Writer) appendOrNext(data []byte) error {
	if err := wr.pageWriter.Append(data); err != nil {
		if err != raw.ErrInsufficientSpace {
			return err
		}

		return wr.next(data)
	}

	return nil
}

func (wr *Writer) next(data []byte) error {
	// flush and start new page
	if err := wr.Flush(); err != nil {
		return err
	}

	wr.pageID++
	// TODO: Test this
	if wr.pageID >= wr.segmentSize {
		return raw.ErrInsufficientSpace
	}

	wr.pageWriter.Init(wr.pageID, wr.timeLineID, record.NewLength(0))
	return wr.pageWriter.Append(data)
}

func (wr *Writer) Flush() error {
	return wr.pageWriter.Write(wr.writer)
}

// TODO: Test this
func (wr *Writer) Initialize(w io.WriterAt) {
	wr.reset(w)

	wr.pageID = 0
	wr.pageWriter.Init(wr.pageID, wr.timeLineID, record.NewLength(0))
}

// TODO: Test this
func (wr *Writer) SeekTail(size int64, rw readWriterAt) error {
	if size == 0 {
		wr.Initialize(rw)
		return nil
	}

	wr.reset(rw)

	wr.pageID = page.NewActivePageIDFromSize(size, wr.pageSize)
	return wr.pageWriter.Read(
		io.NewSectionReader(rw, wr.pageID.Position(wr.pageSize), wr.pageSize))
}

func (wr *Writer) reset(w io.WriterAt) {
	wr.writer = w
	wr.recordBuf.Reset()
}
