package segment

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

// TODO: Should we asynchronously prefetch pages?

type Reader struct {
	pageSize   int64
	pageReader *page.Reader
}

func NewReader(
	pageSize int64,
) *Reader {
	return &Reader{
		pageSize:   pageSize,
		pageReader: page.NewReader(pageSize),
	}
}

func (sr *Reader) seek(pid page.PageID, r io.Seeker) error {
	_, err := r.Seek(pid.Size(sr.pageSize), io.SeekStart)
	return err
}

func (sr *Reader) Iterate(r io.ReadSeeker) iter.Seq2[*record.Record, error] {
	return sr.iterateFrom(0, r)
}

func (sr *Reader) iterateFrom(pid page.PageID, r io.ReadSeeker) iter.Seq2[*record.Record, error] {
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

func (sr *Reader) Reverse(size int64, r io.ReadSeeker) iter.Seq2[*record.Record, error] {
	return sr.reverseFrom(page.NewActivePageIDFromSize(size, sr.pageSize), r)
}

func (sr *Reader) reverseFrom(pid page.PageID, r io.ReadSeeker) iter.Seq2[*record.Record, error] {
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

func (sr *Reader) readForward(pid page.PageID, r io.ReadSeeker) iter.Seq2[iter.Seq2[*record.Record, error], error] {
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

func (sr *Reader) readReverse(pid page.PageID, r io.ReadSeeker) iter.Seq2[iter.Seq2[*record.Record, error], error] {
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
