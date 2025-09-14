package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/raw"
)

type Log struct {
	pageSize  int64
	highWater LogSequenceNumber
	lowWater  LogSequenceNumber
	f         file.File
	recordBuf *bytes.Buffer
	page      *LogPage
}

func (l *Log) HighWater() LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() LogSequenceNumber  { return l.lowWater }

func (l *Log) Open(f file.File) {
	l.f = f
	l.recordBuf = bytes.NewBuffer(nil)
	l.page = NewLogPage(l.pageSize)
}

func (l *Log) IteratePages() iter.Seq2[[]byte, error] {
	_, _ = l.f.Seek(0, 0)
	buffer := make([]byte, l.pageSize)

	return func(yield func([]byte, error) bool) {
		for {
			err := l.readPage(buffer)
			if err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(buffer, nil) {
				return
			}
		}
	}
}

func (l *Log) readPage(buf []byte) error {
	n, err := l.f.Read(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		clear(buf[n:])
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (l *Log) Iterate() iter.Seq2[*LogRecord, error] {
	_, _ = l.f.Seek(0, 0)
	// b := make([]byte, l.pageSize)
	r := bufio.NewReader(l.f)
	return func(yield func(*LogRecord, error) bool) {
		for {
			buffered := r.Buffered()
			fmt.Println(buffered)
			if err := l.validateCRC(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			lr := &LogRecord{}
			err := lr.Read(r)
			if err != nil {
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
			if !yield(lr, nil) || err == io.EOF {
				return
			}
		}
	}
}

func (*Log) validateCRC(r *bufio.Reader) error {
	var c CRC
	if err := c.Read(r); err != nil {
		return err
	}

	lrl := LogRecordLength(0)
	if err := lrl.Read(r); err != nil {
		return err
	}

	if lrl == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(lrl))
	if err != nil {
		return err
	}

	if !c.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}

func (l *Log) reset() {
	// l.buffer.Reset(l.f)
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

func (l *Log) appendPage(lp *LogPage) error {
	n, err := l.f.Write(lp.Data())
	if err != nil {
		return err
	}
	if n < int(l.pageSize) {
		return io.ErrShortWrite
	}
	return nil
}

func (l *Log) Flush(lsn LogSequenceNumber) error {
	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
}
