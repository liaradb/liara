package log

import (
	"bufio"
	"encoding/binary"

	"github.com/cardboardrobots/liaradb/file"
	"github.com/cardboardrobots/liaradb/raw"
)

type Log struct {
	highWater LogSequenceNumber
	lowWater  LogSequenceNumber
	f         file.File
	buffer    *bufio.Writer
}

type LogSequenceNumber uint64

func (l *Log) HighWater() LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() LogSequenceNumber  { return l.lowWater }

func (l *Log) Open(f file.File) {
	l.f = f
	l.buffer = bufio.NewWriter(f)
}

func (l *Log) reset() {
	l.buffer.Reset(l.f)
}

func (l *Log) Append(data []byte) (LogSequenceNumber, error) {
	if err := l.write(data); err != nil {
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
