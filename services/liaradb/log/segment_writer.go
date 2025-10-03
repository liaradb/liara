package log

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

type SegmentWriter struct {
	pageSize    int64
	segmentSize page.PageID
	pageID      page.PageID
	timeLineID  page.TimeLineID
	readWriter  io.ReadWriteSeeker
	recordBuf   *bytes.Buffer
	pageWriter  *PageWriter
}

func NewSegmentWriter(
	pageSize int64,
	segmentSize page.PageID,
	rw io.ReadWriteSeeker,
) *SegmentWriter {
	return &SegmentWriter{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		readWriter:  rw,
		recordBuf:   bytes.NewBuffer(nil),
	}
}

func (sw *SegmentWriter) PageID() page.PageID { return sw.pageID }

func (sw *SegmentWriter) Append(rc *record.Record) error {
	data, err := sw.recordToBytes(rc)
	if err != nil {
		return err
	}

	return sw.append(data)
}

func (sw *SegmentWriter) recordToBytes(rc *record.Record) ([]byte, error) {
	sw.recordBuf.Reset()
	if err := rc.Write(sw.recordBuf); err != nil {
		return nil, err
	}

	return sw.recordBuf.Bytes(), nil
}

func (sw *SegmentWriter) append(data []byte) error {
	rb := record.NewRecordBoundary(data)
	if err := sw.appendOrNext(rb, data); err != nil {
		return err
	}

	return nil
}

func (sw *SegmentWriter) appendOrNext(rb record.RecordBoundary, data []byte) error {
	if err := sw.pageWriter.append(rb, data); err != nil {
		if err != ErrInsufficientSpace {
			return err
		}

		return sw.next(rb, data)
	}

	return nil
}

func (sw *SegmentWriter) next(rb record.RecordBoundary, data []byte) error {
	// flush and start new page
	// TODO: Can we use Write, or do we need Flush?
	if err := sw.pageWriter.Flush(sw.readWriter); err != nil {
		return err
	}

	sw.pageID++
	// TODO: Test this
	if sw.pageID >= sw.segmentSize {
		return ErrInsufficientSpace
	}

	// TODO: Don't replace LogPageWriter
	sw.pageWriter = newPageWriter(sw.pageSize)
	sw.pageWriter.init(sw.pageID, sw.timeLineID, 0)
	return sw.pageWriter.append(rb, data)
}

func (sw *SegmentWriter) Flush() error {
	if err := sw.pageWriter.Flush(sw.readWriter); err != nil {
		return err
	}

	return nil
}

// TODO: Test this
func (sw *SegmentWriter) Initialize() error {
	// TODO: Do we need to seek?
	_, err := sw.readWriter.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	sw.pageID = 0
	// TODO: Don't replace LogPageWriter
	sw.pageWriter = newPageWriter(sw.pageSize)
	sw.pageWriter.init(sw.pageID, sw.timeLineID, 0)

	return nil
}

// TODO: Test this
func (sw *SegmentWriter) SeekTail(size int64) error {
	if size == 0 {
		return sw.Initialize()
	}

	pid := page.NewActivePageIDFromSize(size, sw.pageSize)
	_, err := sw.readWriter.Seek(pid.Size(sw.pageSize), io.SeekStart)
	if err != nil {
		return err
	}

	sw.pageID = pid

	// TODO: initialize or jump to tail of Page
	// Is page initialized?
	// TODO: Don't replace LogPageWriter
	sw.pageWriter = newPageWriter(sw.pageSize)
	sw.pageWriter.init(sw.pageID, sw.timeLineID, 0)

	return sw.pageWriter.SeekTail(sw.readWriter)
}
