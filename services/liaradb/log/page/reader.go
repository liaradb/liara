package page

import (
	"bufio"
	"bytes"
	"container/list"
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/record"
)

type Reader struct {
	data       []byte
	pageReader *bytes.Reader
	pageHeader Header
}

func NewReader(
	pageSize int64,
) *Reader {
	data := make([]byte, pageSize)
	return &Reader{
		data:       data,
		pageReader: bytes.NewReader(data),
	}
}

func (rd *Reader) Header() Header {
	return rd.pageHeader
}

func (rd *Reader) Iterate(r io.Reader) (iter.Seq2[*record.Record, error], error) {
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

// TODO: Change page structure to make reversing easier
func (rd *Reader) Reverse(r io.Reader) (iter.Seq2[*record.Record, error], error) {
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

// TODO: Should we asynchronously prefetch pages?
func (rd *Reader) read(r io.Reader) (*Header, error) {
	if _, err := r.Read(rd.data); err != nil {
		return nil, err
	}

	rd.reset()
	if err := rd.pageHeader.Read(rd.pageReader); err != nil {
		return nil, err
	}

	return &rd.pageHeader, nil
}

func (rd *Reader) reset() {
	rd.pageReader.Reset(rd.data)
}

func (rd *Reader) records() iter.Seq2[*record.Record, error] {
	r := bufio.NewReader(rd.pageReader)
	rb := record.Boundary{}
	return func(yield func(*record.Record, error) bool) {
		for {
			var err error
			// TODO: This reads past the end of the file
			if err = rb.Validate(r); err != nil {
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
