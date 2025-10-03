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

type PageReader struct {
	data       []byte
	pageReader *bytes.Reader
	pageHeader page.Header
}

func NewPageReader(
	pageSize int64,
) *PageReader {
	data := make([]byte, pageSize)
	return &PageReader{
		data:       data,
		pageReader: bytes.NewReader(data),
	}
}

func (pr *PageReader) Header() page.Header {
	return pr.pageHeader
}

func (pr *PageReader) Iterate(r io.Reader) (iter.Seq2[*record.Record, error], error) {
	if _, err := pr.read(r); err != nil {
		return nil, err
	}

	return func(yield func(*record.Record, error) bool) {
		for rc, err := range pr.records() {
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
func (pr *PageReader) Reverse(rd io.Reader) (iter.Seq2[*record.Record, error], error) {
	if _, err := pr.read(rd); err != nil {
		return nil, err
	}

	r := list.New()
	for rc, err := range pr.records() {
		if err != nil {
			return nil, err
		}

		r.PushBack(rc)
	}

	return func(yield func(*record.Record, error) bool) {
		for e := r.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value.(*record.Record), nil) {
				return
			}
		}
	}, nil
}

// TODO: Should we asynchronously prefetch pages?
func (pr *PageReader) read(rd io.Reader) (*page.Header, error) {
	if _, err := rd.Read(pr.data); err != nil {
		return nil, err
	}

	pr.reset()
	if err := pr.pageHeader.Read(pr.pageReader); err != nil {
		return nil, err
	}

	return &pr.pageHeader, nil
}

func (pr *PageReader) reset() {
	pr.pageReader.Reset(pr.data)
}

func (pr *PageReader) records() iter.Seq2[*record.Record, error] {
	r := bufio.NewReader(pr.pageReader)

	return func(yield func(*record.Record, error) bool) {
		for {
			var err error
			// TODO: This reads past the end of the file
			if err = page.ValidateCRC(r); err != nil {
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
