package page

import (
	"container/list"
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/record"
)

type Reader struct {
	page *page.Page[*Header, *page.Item]
}

// TODO: Merge with [storage/record.Page].

func NewReader(pageSize int64) *Reader {
	// TODO: Using three slices/buffers is slow
	return &Reader{
		page: page.NewWithHeader(
			page.Offset(pageSize),
			&Header{},
			page.NewItemByLength),
	}
}

func (rd *Reader) Header() Header { return *rd.page.Header() }

func (rd *Reader) Iterate(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if _, err := rd.read(r); err != nil {
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

// TODO: Change use new page structure to make reversing easier
func (rd *Reader) Reverse(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if _, err := rd.read(r); err != nil {
		return nil, err
	}

	rcs := list.New()
	for rc, err := range rd.records() {
		if err != nil {
			return nil, err
		}

		rcs.PushBack(rc)
	}

	return func(yield func(*record.Record, error) bool) {
		for e := rcs.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value.(*record.Record), nil) {
				return
			}
		}
	}, nil
}

func (rd *Reader) read(r io.ReadSeeker) (*Header, error) {
	if err := rd.page.Read(r); err != nil {
		return nil, err
	}

	return rd.page.Header(), nil
}

func (rd *Reader) records() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for i, err := range rd.page.Items() {
			if err != nil {
				yield(nil, err)
				return
			}

			b := raw.NewBufferFromSlice(i.Value())
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
