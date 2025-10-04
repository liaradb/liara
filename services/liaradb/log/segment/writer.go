package segment

import (
	"bytes"
	"io"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
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
	readWriter  io.ReadWriteSeeker
	recordBuf   *bytes.Buffer
	pageWriter  *page.Writer
}

func NewWriter(
	pageSize int64,
	segmentSize page.PageID,
	rw io.ReadWriteSeeker,
) *Writer {
	return &Writer{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		readWriter:  rw,
		recordBuf:   bytes.NewBuffer(nil),
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

	return wr.recordBuf.Bytes(), nil
}

func (wr *Writer) append(data []byte) error {
	rb := record.NewBoundary(data)
	if err := wr.appendOrNext(rb, data); err != nil {
		return err
	}

	return nil
}

func (wr *Writer) appendOrNext(rb record.Boundary, data []byte) error {
	if err := wr.pageWriter.Append(rb, data); err != nil {
		if err != page.ErrInsufficientSpace {
			return err
		}

		return wr.next(rb, data)
	}

	return nil
}

func (wr *Writer) next(rb record.Boundary, data []byte) error {
	// flush and start new page
	// TODO: Can we use Write, or do we need Flush?
	if err := wr.pageWriter.Flush(wr.readWriter); err != nil {
		return err
	}

	wr.pageID++
	// TODO: Test this
	if wr.pageID >= wr.segmentSize {
		return page.ErrInsufficientSpace
	}

	// TODO: Don't replace LogPageWriter
	wr.pageWriter = page.NewWriter(wr.pageSize)
	wr.pageWriter.Init(wr.pageID, wr.timeLineID, 0)
	return wr.pageWriter.Append(rb, data)
}

func (wr *Writer) Flush() error {
	if err := wr.pageWriter.Flush(wr.readWriter); err != nil {
		return err
	}

	return nil
}

// TODO: Test this
func (wr *Writer) Initialize() error {
	// TODO: Do we need to seek?
	_, err := wr.readWriter.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	wr.pageID = 0
	// TODO: Don't replace LogPageWriter
	wr.pageWriter = page.NewWriter(wr.pageSize)
	wr.pageWriter.Init(wr.pageID, wr.timeLineID, 0)

	return nil
}

// TODO: Test this
func (wr *Writer) SeekTail(size int64) error {
	if size == 0 {
		return wr.Initialize()
	}

	pid := page.NewActivePageIDFromSize(size, wr.pageSize)
	_, err := wr.readWriter.Seek(pid.Size(wr.pageSize), io.SeekStart)
	if err != nil {
		return err
	}

	wr.pageID = pid

	// TODO: initialize or jump to tail of Page
	// Is page initialized?
	// TODO: Don't replace LogPageWriter
	wr.pageWriter = page.NewWriter(wr.pageSize)
	wr.pageWriter.Init(wr.pageID, wr.timeLineID, 0)

	return wr.pageWriter.SeekTail(wr.readWriter)
}
