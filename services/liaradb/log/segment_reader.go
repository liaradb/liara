package log

import (
	"bufio"
	"bytes"
	"container/list"
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type SegmentReader struct {
	size       int64
	reader     io.ReadSeeker
	data       []byte
	pageReader *bytes.Reader
	pageHeader page.PageHeader
}

func NewSegmentReader(
	size int64,
	r io.ReadSeeker,
) *SegmentReader {
	body := size - page.PageHeaderSize
	return &SegmentReader{
		size:   body,
		reader: r,
		data:   make([]byte, body),
	}
}

func (sr *SegmentReader) Seek(pid page.PageID) error {
	_, err := sr.reader.Seek(sr.position(pid, sr.size), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (sr *SegmentReader) position(pid page.PageID, size int64) int64 {
	return int64(pid) * (size + page.PageHeaderSize)
}

func (sr *SegmentReader) Iterate() iter.Seq2[*record.Record, error] {
	return sr.IterateFrom(0)
}

// TODO: Test this
func (sr *SegmentReader) IterateFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		if err := sr.Seek(pid); err != nil {
			yield(nil, err)
			return
		}

		for {
			if _, err := sr.Read(); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			for rc, err := range sr.Records() {
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
func (sr *SegmentReader) Reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		q := list.New()

		for rc, err := range sr.Iterate() {
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
func (sr *SegmentReader) Read() (*page.PageHeader, error) {
	if err := sr.pageHeader.Read(sr.reader); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := sr.reader.Read(sr.data); err != nil {
		return nil, err
	}

	sr.initReader()

	return &sr.pageHeader, nil
}

func (sr *SegmentReader) initReader() {
	if sr.pageReader == nil {
		sr.pageReader = bytes.NewReader(sr.data)
	} else {
		sr.pageReader.Reset(sr.data)
	}
}

func (sr *SegmentReader) Records() iter.Seq2[*record.Record, error] {
	r := bufio.NewReader(sr.pageReader)

	return func(yield func(*record.Record, error) bool) {
		for {
			var err error
			if err = sr.validateCRC(r); err != nil {
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

func (*SegmentReader) validateCRC(r *bufio.Reader) error {
	var c page.CRC
	if err := c.Read(r); err != nil {
		return err
	}

	rl := page.RecordLength(0)
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
