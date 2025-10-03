package log

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type SegmentReader struct {
	pageSize int64
	reader   io.ReadSeeker
	pReader  *PageReader
}

func NewSegmentReader(
	pageSize int64,
	r io.ReadSeeker,
) *SegmentReader {
	return &SegmentReader{
		pageSize: pageSize,
		reader:   r,
		pReader:  NewPageReader(pageSize),
	}
}

func (sr *SegmentReader) seek(pid page.PageID) error {
	_, err := sr.reader.Seek(pid.Size(sr.pageSize), io.SeekStart)
	return err
}

func (sr *SegmentReader) Iterate() iter.Seq2[*record.Record, error] {
	return sr.iterateFrom(0)
}

func (sr *SegmentReader) iterateFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
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
	return sr.reverseFrom(page.NewActivePageIDFromSize(size, sr.pageSize))
}

func (sr *SegmentReader) reverseFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
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
			if err := sr.seek(pid); err != nil {
				yield(nil, err)
				return
			}

			it, err := sr.pReader.Iterate(sr.reader)
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
			if err := sr.seek(pid - i); err != nil {
				yield(nil, err)
				return
			}

			it, err := sr.pReader.Reverse(sr.reader)
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
