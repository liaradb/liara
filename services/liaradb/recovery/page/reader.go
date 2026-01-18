package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/mempage"
	"github.com/liaradb/liaradb/recovery/record"
)

type Reader struct {
	page Page[*mempage.Item]
}

func NewReader(pageSize int64) *Reader {
	return &Reader{
		page: mempage.NewWithHeader(
			page.Offset(pageSize),
			&Header{},
			mempage.NewItemByLength),
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

func (rd *Reader) Reverse(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if _, err := rd.read(r); err != nil {
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

func (rd *Reader) reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for i, err := range rd.page.ItemsReverse() {
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
