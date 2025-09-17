package log

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

type LogReader struct {
	size       int64
	reader     io.ReadSeeker
	data       []byte
	pageReader *bytes.Reader
	page       LogPageHeader
}

func NewLogReader(
	size int64,
	r io.ReadSeeker,
) *LogReader {
	body := size - pageHeaderSize
	return &LogReader{
		size:   body,
		reader: r,
		data:   make([]byte, body),
	}
}

func (l *LogReader) Seek(pid LogPageID) error {
	_, err := l.reader.Seek(l.position(pid, l.size), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (l *LogReader) position(pid LogPageID, size int64) int64 {
	return int64(pid) * (size + pageHeaderSize)
}

func (l *LogReader) Iterate() iter.Seq2[*LogRecord, error] {
	return l.IterateFrom(0)
}

func (l *LogReader) IterateFrom(pid LogPageID) iter.Seq2[*LogRecord, error] {
	return func(yield func(*LogRecord, error) bool) {
		if err := l.Seek(pid); err != nil {
			yield(nil, err)
			return
		}

		for {
			if _, err := l.Read(); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			for lr, err := range l.Records() {
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

func (l *LogReader) Read() (*LogPageHeader, error) {
	if err := l.page.Read(l.reader); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := l.reader.Read(l.data); err != nil {
		return nil, err
	}

	l.initReader()

	return &l.page, nil
}

func (l *LogReader) initReader() {
	if l.pageReader == nil {
		l.pageReader = bytes.NewReader(l.data)
	} else {
		l.pageReader.Reset(l.data)
	}
}

func (l *LogReader) Records() iter.Seq2[*LogRecord, error] {
	r := bufio.NewReader(l.pageReader)
	lr := &LogRecord{}

	return func(yield func(*LogRecord, error) bool) {
		for {
			var err error
			if err = l.validateCRC(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			// TODO: Use a buffer
			if err := lr.Read(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(lr, nil) {
				return
			}
		}
	}
}

func (*LogReader) validateCRC(r *bufio.Reader) error {
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
