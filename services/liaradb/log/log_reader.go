package log

import (
	"bufio"
	"bytes"
	"container/list"
	"io"
	"iter"
)

type LogReader struct {
	size       int64
	reader     io.ReadSeeker
	data       []byte
	pageReader *bytes.Reader
	header     PageHeader
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

func (l *LogReader) Seek(pid PageID) error {
	_, err := l.reader.Seek(l.position(pid, l.size), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (l *LogReader) position(pid PageID, size int64) int64 {
	return int64(pid) * (size + pageHeaderSize)
}

func (l *LogReader) Iterate() iter.Seq2[*Record, error] {
	return l.IterateFrom(0)
}

// TODO: Test this
func (l *LogReader) IterateFrom(pid PageID) iter.Seq2[*Record, error] {
	return func(yield func(*Record, error) bool) {
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

			for rc, err := range l.Records() {
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(rc, nil) {
					return
				}
			}
		}
	}
}

// TODO: Change page structure to make reversing easier
func (l *LogReader) Reverse() iter.Seq2[*Record, error] {
	return func(yield func(*Record, error) bool) {
		q := list.New()

		for rc, err := range l.Iterate() {
			if err != nil {
				yield(nil, err)
				return
			}

			q.PushBack(rc)
		}

		for {
			e := q.Back()
			if e == nil {
				return
			}

			v := e.Value.(*Record)
			q.Remove(e)
			if !yield(v, nil) {
				return
			}
		}
	}
}

// TODO: Should we asynchronously prefetch pages?
func (l *LogReader) Read() (*PageHeader, error) {
	if err := l.header.Read(l.reader); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := l.reader.Read(l.data); err != nil {
		return nil, err
	}

	l.initReader()

	return &l.header, nil
}

func (l *LogReader) initReader() {
	if l.pageReader == nil {
		l.pageReader = bytes.NewReader(l.data)
	} else {
		l.pageReader.Reset(l.data)
	}
}

func (l *LogReader) Records() iter.Seq2[*Record, error] {
	r := bufio.NewReader(l.pageReader)

	return func(yield func(*Record, error) bool) {
		for {
			var err error
			if err = l.validateCRC(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			// TODO: Should we create a new record each time?
			rc := &Record{}

			// TODO: Use a buffer
			if err := rc.Read(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(rc, nil) {
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

	rl := RecordLength(0)
	if err := rl.Read(r); err != nil {
		return err
	}

	if rl == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(rl))
	if err != nil {
		return err
	}

	if !c.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}
