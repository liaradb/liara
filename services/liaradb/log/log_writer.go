package log

import (
	"bytes"
	"io"
)

const (
	BlockSize   = 1024
	SegmentSize = 1024
)

type LogWriter struct {
	pageSize   int64
	pageIndex  LogPageID
	timeLineID TimeLineID
	highWater  LogSequenceNumber
	lowWater   LogSequenceNumber
	writer     io.WriteSeeker
	recordBuf  *bytes.Buffer
	page       *LogPageWriter
}

func NewLogWriter(pageSize int64, w io.WriteSeeker) *LogWriter {
	return &LogWriter{
		pageSize:  pageSize,
		writer:    w,
		recordBuf: bytes.NewBuffer(nil),
		page:      newLogPageWriter(pageSize),
	}
}

func (l *LogWriter) PageIndex() LogPageID         { return l.pageIndex }
func (l *LogWriter) HighWater() LogSequenceNumber { return l.highWater }
func (l *LogWriter) LowWater() LogSequenceNumber  { return l.lowWater }

func (l *LogWriter) Append(lr *Record) (LogSequenceNumber, error) {
	data, err := l.recordToBytes(lr)
	if err != nil {
		return 0, err
	}

	return l.append(data)
}

func (l *LogWriter) recordToBytes(lr *Record) ([]byte, error) {
	l.recordBuf.Reset()
	if err := lr.Write(l.recordBuf); err != nil {
		return nil, err
	}

	return l.recordBuf.Bytes(), nil
}

func (l *LogWriter) append(data []byte) (LogSequenceNumber, error) {
	crc := NewCRC(data)
	if err := crc.Write(l.writer); err != nil {
		return 0, err
	}

	if err := l.appendOrNext(crc, data); err != nil {
		return 0, err
	}

	l.highWater++
	return l.highWater, nil
}

func (l *LogWriter) appendOrNext(crc CRC, data []byte) error {
	if err := l.page.append(crc, data); err != nil {
		if err != ErrInsufficientSpace {
			return err
		}

		return l.next(crc, data)
	}

	return nil
}

func (l *LogWriter) next(crc CRC, data []byte) error {
	// flush and start new page
	// TODO: Can we use Write, or do we need Flush?
	if err := l.page.Flush(l.writer); err != nil {
		return err
	}

	l.pageIndex++
	// TODO: Don't replace LogPageWriter
	l.page = newLogPageWriter(l.pageSize)
	l.page.init(l.pageIndex, l.timeLineID, 0)
	return l.page.append(crc, data)
}

func (l *LogWriter) Flush(lsn LogSequenceNumber) error {
	if err := l.page.Flush(l.writer); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
}
