package log

import (
	"bufio"
	"bytes"
	"encoding/binary"
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
		// var c CRC
		// if err := c.Read(l.f); err != nil {
		// 	yield(nil, err)
		// 	return
		// }

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

	// crc := NewCRC(l.rb.Bytes())
	// if err := crc.Write(l.rb); err != nil {
	// 	return 0, err
	// }

	return l.append(l.rb.Bytes())
}

func (l *Log) append(data []byte) (LogSequenceNumber, error) {
	if _, err := l.f.Write(data); err != nil {
		l.reset()
		return 0, err
	}

	l.highWater++
	return l.highWater, nil
}

func (l *Log) write(data []byte) error {
	if err := l.writeSize(data); err != nil {
		return err
	}

	return l.writeData(data)
}

func (l *Log) writeSize(data []byte) error {
	return binary.Write(l.buffer, binary.BigEndian, uint32(len(data)))
}

func (l *Log) writeData(data []byte) error {
	n, err := l.buffer.Write(data)
	if n != len(data) {
		return raw.ErrOverflow
	}

	return err
}

func (l *Log) Flush(lsn LogSequenceNumber) error {
	if err := l.buffer.Flush(); err != nil {
		return err
	}

	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
}
