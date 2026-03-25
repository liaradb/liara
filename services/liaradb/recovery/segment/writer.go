package segment

import (
	"bytes"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	BlockSize   = 1024
	SegmentSize = 1024
)

type Writer struct {
	pageSize    int64
	segmentSize action.PageID
	pageID      action.PageID
	timeLineID  action.TimeLineID
	writer      io.WriterAt
	recordBuf   *bytes.Buffer
	pageWriter  *page.Page
}

type readWriterAt interface {
	io.WriterAt
	io.ReaderAt
}

func NewWriter(
	pageSize int64,
	segmentSize action.PageID,
) *Writer {
	return &Writer{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		recordBuf:   bytes.NewBuffer(nil),
		pageWriter:  page.New(pageSize),
	}
}

func (wr *Writer) PageID() action.PageID { return wr.pageID }

func (wr *Writer) Append(rc *record.Record) error {
	data, err := wr.recordToBytes(rc)
	if err != nil {
		return err
	}

	return wr.appendOrNext(data)
}

func (wr *Writer) recordToBytes(rc *record.Record) ([]byte, error) {
	wr.recordBuf.Reset()
	if err := rc.Write(wr.recordBuf); err != nil {
		return nil, err
	}

	// We don't need to clone, as the data is copied
	return wr.recordBuf.Bytes(), nil
}

func (wr *Writer) appendOrNext(data []byte) error {
	if ok := wr.pageWriter.Append(data); ok {
		return nil
	}

	return wr.next(data)
}

func (wr *Writer) next(data []byte) error {
	// flush and start new page
	if err := wr.Flush(); err != nil {
		return err
	}

	wr.pageID++
	if wr.pageID >= wr.segmentSize {
		return raw.ErrInsufficientSpace
	}

	// TODO: Verify that record can fit at all before initializing
	wr.pageWriter.Init(wr.pageID, wr.timeLineID, record.NewLength(0))
	if ok := wr.pageWriter.Append(data); !ok {
		return raw.ErrInsufficientSpace
	}

	return nil
}

func (wr *Writer) Flush() error {
	return wr.pageWriter.Write(wr.writer)
}

func (wr *Writer) SeekTail(size int64, rw readWriterAt) error {
	if size == 0 {
		wr.initialize(rw)
		return nil
	}

	wr.reset(rw)

	wr.pageID = action.NewActivePageIDFromSize(size, wr.pageSize)
	return wr.pageWriter.Read(
		io.NewSectionReader(rw, wr.pageID.Position(wr.pageSize), wr.pageSize))
}

func (wr *Writer) initialize(w io.WriterAt) {
	wr.reset(w)

	wr.pageID = 0
	wr.pageWriter.Init(wr.pageID, wr.timeLineID, record.NewLength(0))
}

func (wr *Writer) reset(w io.WriterAt) {
	wr.writer = w
	wr.recordBuf.Reset()
}
