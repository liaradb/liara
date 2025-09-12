package log

import (
	"bufio"
	"bytes"
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/raw"
)

type Log struct {
	pageSize  int64
	highWater LogSequenceNumber
	lowWater  LogSequenceNumber
	f         file.File
	buffer    *bufio.Writer
	rb        *bytes.Buffer
}

func (l *Log) HighWater() LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() LogSequenceNumber  { return l.lowWater }

func (l *Log) Open(f file.File) {
	l.f = f
	l.buffer = bufio.NewWriter(f)
	l.rb = bytes.NewBuffer(nil)
}

func (l *Log) Iterate() iter.Seq2[*LogRecord, error] {
	_, _ = l.f.Seek(0, 0)
	// b := make([]byte, l.pageSize)
	return func(yield func(*LogRecord, error) bool) {
		var c CRC
		if err := c.Read(l.f); err != nil {
			yield(nil, err)
			return
		}

		lrl := LogRecordLength(0)
		if err := lrl.Read(l.f); err != nil {
			yield(nil, err)
			return
		}

		lr := &LogRecord{}
		if err := lr.Read(l.f); err != nil {
			yield(nil, err)
			return
		}

		// ld := &LogData{}
		// err := ld.Read(l.f)
		// if err != nil {
		// 	yield(nil, err)
		// 	return
		// }
		// if n < int(l.pageSize) {
		// 	for i := n; i < int(l.pageSize); i++ {
		// 		b[i] = 0
		// 	}
		// }
		if !yield(lr, nil) {
			return
		}
	}
}

func (l *Log) reset() {
	l.buffer.Reset(l.f)
}

func (l *Log) Append(lr *LogRecord) (LogSequenceNumber, error) {
	l.rb.Reset()
	if err := lr.Write(l.rb); err != nil {
		return 0, err
	}

	return l.append(l.rb.Bytes())
}

func (l *Log) append(data []byte) (LogSequenceNumber, error) {
	crc := NewCRC(l.rb.Bytes())
	if err := crc.Write(l.f); err != nil {
		return 0, err
	}

	if err := NewLogRecordLength(data).Write(l.f); err != nil {
		l.reset()
		return 0, err
	}

	// TODO: Do we need to verify write lengths?
	if n, err := l.f.Write(data); err != nil {
		l.reset()
		return 0, err
	} else if n != len(data) {
		return 0, raw.ErrOverflow
	}

	l.highWater++
	return l.highWater, nil
}

func (l *Log) Flush(lsn LogSequenceNumber) error {
	if err := l.buffer.Flush(); err != nil {
		return err
	}

	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
}
