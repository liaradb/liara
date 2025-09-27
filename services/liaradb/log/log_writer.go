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

type LogWriter struct {
	pageSize    int64
	segmentSize page.PageID
	pageID      page.PageID
	timeLineID  page.TimeLineID
	highWater   record.LogSequenceNumber
	lowWater    record.LogSequenceNumber
	writer      io.WriteSeeker
	recordBuf   *bytes.Buffer
	pageWriter  *PageWriter
}

func NewLogWriter(
	pageSize int64,
	segmentSize page.PageID,
	w io.WriteSeeker,
) *LogWriter {
	return &LogWriter{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		writer:      w,
		recordBuf:   bytes.NewBuffer(nil),
		pageWriter:  newPageWriter(pageSize),
	}
}

func (lw *LogWriter) PageID() page.PageID                 { return lw.pageID }
func (lw *LogWriter) HighWater() record.LogSequenceNumber { return lw.highWater }
func (lw *LogWriter) LowWater() record.LogSequenceNumber  { return lw.lowWater }

func (lw *LogWriter) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	data, err := lw.recordToBytes(rc)
	if err != nil {
		return 0, err
	}

	return lw.append(data)
}

func (lw *LogWriter) recordToBytes(rc *record.Record) ([]byte, error) {
	lw.recordBuf.Reset()
	if err := rc.Write(lw.recordBuf); err != nil {
		return nil, err
	}

	return lw.recordBuf.Bytes(), nil
}

func (lw *LogWriter) append(data []byte) (record.LogSequenceNumber, error) {
	crc := page.NewCRC(data)
	if err := crc.Write(lw.writer); err != nil {
		return 0, err
	}

	if err := lw.appendOrNext(crc, data); err != nil {
		if err == ErrInsufficientSpace {
			// TODO: Fix this
			return lw.highWater + 1, err
		}
		return 0, err
	}

	lw.highWater++
	return lw.highWater, nil
}

func (lw *LogWriter) appendOrNext(crc page.CRC, data []byte) error {
	if err := lw.pageWriter.append(crc, data); err != nil {
		if err != ErrInsufficientSpace {
			return err
		}

		return lw.next(crc, data)
	}

	return nil
}

func (lw *LogWriter) next(crc page.CRC, data []byte) error {
	// flush and start new page
	// TODO: Can we use Write, or do we need Flush?
	if err := lw.pageWriter.Flush(lw.writer); err != nil {
		return err
	}

	lw.pageID++
	if lw.pageID > lw.segmentSize {
		return ErrInsufficientSpace
	}

	// TODO: Don't replace LogPageWriter
	lw.pageWriter = newPageWriter(lw.pageSize)
	lw.pageWriter.init(lw.pageID, lw.timeLineID, 0)
	return lw.pageWriter.append(crc, data)
}

func (lw *LogWriter) Flush(lsn record.LogSequenceNumber) error {
	if err := lw.pageWriter.Flush(lw.writer); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, lw.highWater)
	lw.lowWater = lsn
	return nil
}

// TODO: Test this
func (lw *LogWriter) SeekTail(size int64) error {
	_, err := lw.writer.Seek(size, io.SeekStart)
	if err != nil {
		return err
	}

	if err := lw.pageWriter.SeekTail(); err != nil {
		return nil
	}

	// TODO: Jump to tail of Page
	lw.pageID = lw.sizeToPageID(size)

	return err
}

// TODO: This is also defined on SegmentReader, but subtracted by 1
func (lw *LogWriter) sizeToPageID(size int64) page.PageID {
	pid := size / lw.pageSize
	if size%lw.pageSize != 0 {
		pid++
	}
	return page.PageID(pid)
}
