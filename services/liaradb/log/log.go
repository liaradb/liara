package log

import (
	"bytes"
	"io"
	"iter"

	"github.com/liaradb/liaradb/file"
)

const (
	BlockSize        = 1024
	SegmentSize      = 1024
	PageHeaderSize   = 4 + 8 + 4
	RecordHeaderSize = 4 + 4
)

type Log struct {
	pageSize   int64
	pageIndex  LogPageID
	timeLineID TimeLineID
	highWater  LogSequenceNumber
	lowWater   LogSequenceNumber
	f          file.File
	recordBuf  *bytes.Buffer
	page       *LogPage
}

func (l *Log) PageIndex() LogPageID         { return l.pageIndex }
func (l *Log) HighWater() LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() LogSequenceNumber  { return l.lowWater }

func (l *Log) Open(f file.File) {
	l.f = f
	l.recordBuf = bytes.NewBuffer(nil)
	l.page = newLogPage(l.pageSize)
}

func (l *Log) Iterate() iter.Seq2[*LogRecord, error] {
	return l.IterateFrom(0)
}

func (l *Log) IterateFrom(pid LogPageID) iter.Seq2[*LogRecord, error] {
	lp := newLogPage(l.pageSize)
	// TODO: What is the correct TimeLineID?
	lp.init(pid, l.timeLineID)
	if err := lp.Seek(l.f); err != nil {
		return errorIterator[*LogRecord](err)
	}

	return func(yield func(*LogRecord, error) bool) {
		for {
			if err := lp.Read(l.f); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			for lr, err := range lp.Records() {
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(lr, nil) {
					return
				}
			}
		}
	}
}

func (l *Log) Append(lr *LogRecord) (LogSequenceNumber, error) {
	data, err := l.recordToBytes(lr)
	if err != nil {
		return 0, err
	}

	return l.append(data)
}

func (l *Log) recordToBytes(lr *LogRecord) ([]byte, error) {
	l.recordBuf.Reset()
	if err := lr.Write(l.recordBuf); err != nil {
		return nil, err
	}

	return l.recordBuf.Bytes(), nil
}

func (l *Log) append(data []byte) (LogSequenceNumber, error) {
	crc := NewCRC(data)
	if err := crc.Write(l.f); err != nil {
		return 0, err
	}

	if err := l.AppendOrNext(crc, data); err != nil {
		return 0, err
	}

	l.highWater++
	return l.highWater, nil
}

func (l *Log) AppendOrNext(crc CRC, data []byte) error {
	err := l.page.append(crc, data)
	if err == nil {
		return nil
	}

	if err == ErrInsufficientSpace {
		// flush and start new page
		// TODO: Can we use Write, or do we need Flush?
		if err := l.page.Flush(l.f); err != nil {
			return err
		}

		l.pageIndex++
		l.page = newLogPage(l.pageSize)
		l.page.init(l.pageIndex, l.timeLineID)
		return l.page.append(crc, data)
	}

	return err
}

func (l *Log) Flush(lsn LogSequenceNumber) error {
	if err := l.page.Flush(l.f); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
}

func errorIterator[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var v T
		yield(v, err)
	}
}
