package log

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type SegmentReader struct {
	pageSize   int64
	bodySize   int64
	reader     io.ReadSeeker
	data       []byte
	pageHeader page.Header
	pReader    *PageReader
}

func NewSegmentReader(
	pageSize int64,
	r io.ReadSeeker,
) *SegmentReader {
	sr := &SegmentReader{
		pageSize: pageSize,
		reader:   r,
	}
	body := pageSize - int64(sr.pageHeader.Size())
	sr.bodySize = body
	sr.data = make([]byte, body)
	sr.pReader = NewPageReader(pageSize, r)
	return sr
}

func (sr *SegmentReader) Seek(pid page.PageID) error {
	_, err := sr.reader.Seek(pid.Size(sr.pageSize), io.SeekStart)
	return err
}

func (sr *SegmentReader) Iterate() iter.Seq2[*record.Record, error] {
	return sr.IterateFrom(0)
}

// TODO: Test this
func (sr *SegmentReader) IterateFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for it, err := range sr.readForward(pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range it {
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
		for it, err := range sr.readReverse(pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range it {
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

func (sr *SegmentReader) readForward(pid page.PageID) iter.Seq2[iter.Seq2[*record.Record, error], error] {
	return func(yield func(iter.Seq2[*record.Record, error], error) bool) {
		for {
			if err := sr.Seek(pid); err != nil {
				yield(nil, err)
				return
			}

			it, err := sr.pReader.Iterate()
			if err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(it, nil) {
				return
			}
			pid++
		}
	}
}

func (sr *SegmentReader) readReverse(pid page.PageID) iter.Seq2[iter.Seq2[*record.Record, error], error] {
	return func(yield func(iter.Seq2[*record.Record, error], error) bool) {
		for i := range pid + 1 {
			if err := sr.Seek(pid - i); err != nil {
				yield(nil, err)
				return
			}

			it, err := sr.pReader.Reverse()
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(it, nil) {
				return
			}
		}
	}
}
