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
	pageSize    int64
	bodySize    int64
	segmentSize page.PageID
	reader      io.ReadSeeker
	data        []byte
	pageReader  *bytes.Reader
	pageHeader  page.PageHeader
}

func NewSegmentReader(
	pageSize int64,
	segmentSize page.PageID,
	r io.ReadSeeker,
) *SegmentReader {
	body := pageSize - page.PageHeaderSize
	return &SegmentReader{
		pageSize:    pageSize,
		bodySize:    body,
		segmentSize: segmentSize,
		reader:      r,
		data:        make([]byte, body),
	}
}

func (sr *SegmentReader) Seek(pid page.PageID) error {
	_, err := sr.reader.Seek(sr.position(pid), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (sr *SegmentReader) position(pid page.PageID) int64 {
	return int64(pid) * sr.pageSize
}

func (sr *SegmentReader) Iterate() iter.Seq2[*record.Record, error] {
	return sr.IterateFrom(0)
}

// TODO: Test this
func (sr *SegmentReader) IterateFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for err := range sr.readForward(pid) {
			if err != nil {
				yield(nil, err)
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

func (sr *SegmentReader) Reverse(size int64) iter.Seq2[*record.Record, error] {
	return sr.ReverseFrom(page.NewActivePageIDFromSize(size, sr.pageSize))
}

// TODO: Change page structure to make reversing easier
// TODO: Test this
func (sr *SegmentReader) ReverseFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for err := range sr.readReverse(pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			r := list.New()
			for rc, err := range sr.Records() {
				if err != nil {
					yield(nil, err)
					return
				}

				r.PushBack(rc)
			}

			for e := r.Back(); e != nil; e = e.Prev() {
				if !yield(e.Value.(*record.Record), nil) {
					return
				}
			}
		}
	}
}

func (sr *SegmentReader) readForward(pid page.PageID) iter.Seq[error] {
	return func(yield func(error) bool) {
		if err := sr.Seek(pid); err != nil {
			yield(err)
			return
		}

		for {
			if _, err := sr.Read(); err != nil {
				if err != io.EOF {
					yield(err)
				}
				return
			}

			if !yield(nil) {
				return
			}
		}
	}
}

func (sr *SegmentReader) readReverse(pid page.PageID) iter.Seq[error] {
	return func(yield func(error) bool) {
		for i := range pid + 1 {
			if _, err := sr.ReadAt(pid - i); err != nil {
				if err != io.EOF {
					yield(err)
				}
				return
			}

			if !yield(nil) {
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

func (sr *SegmentReader) ReadAt(pid page.PageID) (*page.PageHeader, error) {
	if err := sr.Seek(pid); err != nil {
		return nil, err
	}

	return sr.Read()
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
