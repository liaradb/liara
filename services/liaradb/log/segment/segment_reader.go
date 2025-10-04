package segment

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type SegmentReader struct {
	pageSize   int64
	pageReader *page.Reader
}

func NewSegmentReader(
	pageSize int64,
) *SegmentReader {
	return &SegmentReader{
		pageSize:   pageSize,
		pageReader: page.NewReader(pageSize),
	}
}

func (sr *SegmentReader) seek(pid page.PageID, r io.Seeker) error {
	_, err := r.Seek(pid.Size(sr.pageSize), io.SeekStart)
	return err
}

func (sr *SegmentReader) Iterate(r io.ReadSeeker) iter.Seq2[*record.Record, error] {
	return sr.iterateFrom(0, r)
}

func (sr *SegmentReader) iterateFrom(pid page.PageID, r io.ReadSeeker) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for it, err := range sr.readForward(pid, r) {
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

func (sr *SegmentReader) Reverse(size int64, r io.ReadSeeker) iter.Seq2[*record.Record, error] {
	return sr.reverseFrom(page.NewActivePageIDFromSize(size, sr.pageSize), r)
}

func (sr *SegmentReader) reverseFrom(pid page.PageID, r io.ReadSeeker) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for it, err := range sr.readReverse(pid, r) {
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

func (sr *SegmentReader) readForward(pid page.PageID, r io.ReadSeeker) iter.Seq2[iter.Seq2[*record.Record, error], error] {
	return func(yield func(iter.Seq2[*record.Record, error], error) bool) {
		for {
			if err := sr.seek(pid, r); err != nil {
				yield(nil, err)
				return
			}

			it, err := sr.pageReader.Iterate(r)
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

func (sr *SegmentReader) readReverse(pid page.PageID, r io.ReadSeeker) iter.Seq2[iter.Seq2[*record.Record, error], error] {
	return func(yield func(iter.Seq2[*record.Record, error], error) bool) {
		for i := range pid + 1 {
			if err := sr.seek(pid-i, r); err != nil {
				yield(nil, err)
				return
			}

			it, err := sr.pageReader.Reverse(r)
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
