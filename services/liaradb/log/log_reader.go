package log

import (
	"bufio"
	"bytes"
	"container/list"
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type LogReader struct {
	size       int64
	reader     io.ReadSeeker
	sl         *segment.SegmentList
	data       []byte
	pageReader *bytes.Reader
	pageHeader record.PageHeader
}

func NewLogReader(
	size int64,
	sl *segment.SegmentList,
	r io.ReadSeeker,
) *LogReader {
	body := size - record.PageHeaderSize
	return &LogReader{
		size:   body,
		sl:     sl,
		reader: r,
		data:   make([]byte, body),
	}
}

func (lr *LogReader) Seek(pid record.PageID) error {
	_, err := lr.reader.Seek(lr.position(pid, lr.size), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (lr *LogReader) position(pid record.PageID, size int64) int64 {
	return int64(pid) * (size + record.PageHeaderSize)
}

func (lr *LogReader) Iterate() iter.Seq2[*record.Record, error] {
	return lr.IterateFrom(0)
}

// TODO: Test this
func (lr *LogReader) IterateFrom(pid record.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		if err := lr.Seek(pid); err != nil {
			yield(nil, err)
			return
		}

		for {
			if _, err := lr.Read(); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			for rc, err := range lr.Records() {
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
func (lr *LogReader) Reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		q := list.New()

		for rc, err := range lr.Iterate() {
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

			v := e.Value.(*record.Record)
			q.Remove(e)
			if !yield(v, nil) {
				return
			}
		}
	}
}

// TODO: Should we asynchronously prefetch pages?
func (lr *LogReader) Read() (*record.PageHeader, error) {
	if err := lr.pageHeader.Read(lr.reader); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := lr.reader.Read(lr.data); err != nil {
		return nil, err
	}

	lr.initReader()

	return &lr.pageHeader, nil
}

func (lr *LogReader) initReader() {
	if lr.pageReader == nil {
		lr.pageReader = bytes.NewReader(lr.data)
	} else {
		lr.pageReader.Reset(lr.data)
	}
}

func (lr *LogReader) Records() iter.Seq2[*record.Record, error] {
	r := bufio.NewReader(lr.pageReader)

	return func(yield func(*record.Record, error) bool) {
		for {
			var err error
			if err = lr.validateCRC(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			// TODO: Should we create a new record each time?
			rc := &record.Record{}

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
	var c record.CRC
	if err := c.Read(r); err != nil {
		return err
	}

	rl := record.RecordLength(0)
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
