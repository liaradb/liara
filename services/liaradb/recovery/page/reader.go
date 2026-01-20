package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

type Reader struct {
	page Page
}

func NewReader(page Page) *Reader {
	return &Reader{
		page: page,
	}
}

func (h *Reader) ID() action.PageID              { return h.page.ID() }
func (h *Reader) TimeLineID() action.TimeLineID  { return h.page.TimeLineID() }
func (h *Reader) LengthRemaining() record.Length { return h.page.LengthRemaining() }

func (rd *Reader) Iterate(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if err := rd.read(r); err != nil {
		return nil, err
	}

	return func(yield func(*record.Record, error) bool) {
		for rc, err := range rd.records() {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(rc, nil) {
				return
			}
		}
	}, nil
}

func (rd *Reader) Reverse(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if err := rd.read(r); err != nil {
		return nil, err
	}

	return func(yield func(*record.Record, error) bool) {
		for rc, err := range rd.reverse() {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(rc, nil) {
				return
			}
		}
	}, nil
}

func (rd *Reader) read(r io.ReadSeeker) error {
	return rd.page.Read(r)
}

func (rd *Reader) records() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for i := range rd.page.Items() {
			b := raw.NewBufferFromSlice(i)
			rc := &record.Record{}
			if err := rc.Read(b); err != nil {
				yield(nil, err)
				return
			}

			if !yield(rc, nil) {
				return
			}
		}
	}
}

func (rd *Reader) reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for i := range rd.page.ItemsReverse() {
			b := raw.NewBufferFromSlice(i)
			rc := &record.Record{}
			if err := rc.Read(b); err != nil {
				yield(nil, err)
				return
			}

			if !yield(rc, nil) {
				return
			}
		}
	}
}
