package segment

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Reader struct {
	pageSize   int64
	pageReader *page.Page
}

func NewReader(
	pageSize int64,
) *Reader {
	return &Reader{
		pageSize:   pageSize,
		pageReader: page.New(pageSize),
	}
}

func (sr *Reader) position(pid action.PageID) int64 {
	return pid.Position(sr.pageSize)
}

func (sr *Reader) Iterate(r io.ReaderAt) iter.Seq2[*record.Record, error] {
	return sr.iterateFrom(0, r)
}

func (sr *Reader) iterateFrom(pid action.PageID, r io.ReaderAt) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for it, err := range sr.readForward(pid, r) {
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range it {
				if !yield(rc, err) || err != nil {
					return
				}
			}
		}
	}
}

func (sr *Reader) Reverse(size int64, r io.ReaderAt) iter.Seq2[*record.Record, error] {
	return sr.reverseFrom(action.NewActivePageIDFromSize(size, sr.pageSize), r)
}

func (sr *Reader) reverseFrom(pid action.PageID, r io.ReaderAt) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for it, err := range sr.readReverse(pid, r) {
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range it {
				if !yield(rc, nil) || err != nil {
					return
				}
			}
		}
	}
}

func (sr *Reader) readForward(pid action.PageID, r io.ReaderAt) iter.Seq2[iter.Seq2[*record.Record, error], error] {
	return func(yield func(iter.Seq2[*record.Record, error], error) bool) {
		for {
			sec := io.NewSectionReader(r, sr.position(pid), sr.pageSize)
			it, err := sr.pageReader.Iterate(sec)
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

func (sr *Reader) readReverse(pid action.PageID, r io.ReaderAt) iter.Seq2[iter.Seq2[*record.Record, error], error] {
	return func(yield func(iter.Seq2[*record.Record, error], error) bool) {
		for i := range pid + 1 {
			sec := io.NewSectionReader(r, sr.position(pid-i), sr.pageSize)
			it, err := sr.pageReader.Reverse(sec)
			if err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(it, nil) {
				return
			}
		}
	}
}
