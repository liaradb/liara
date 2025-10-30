package page

import (
	"bufio"
	"bytes"
	"container/list"
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Reader struct {
	data       []byte
	pageReader *bytes.Reader
	reader     *bufio.Reader
	pageHeader Header
}

// TODO: Merge with [storage/record.Page].

func NewReader(pageSize int64) *Reader {
	// TODO: Using three slices/buffers is slow
	data := make([]byte, pageSize)
	pageReader := bytes.NewReader(data)
	return &Reader{
		data:       data,
		pageReader: pageReader,
		reader:     bufio.NewReaderSize(nil, int(pageSize)),
	}
}

func (rd *Reader) Header() Header { return rd.pageHeader }

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

func (rd *Reader) read(r io.Reader) (*Header, error) {
	if _, err := r.Read(rd.data); err != nil {
		return nil, err
	}

	rd.pageReader.Reset(rd.data)
	if err := rd.pageHeader.Read(rd.pageReader); err != nil {
		return nil, err
	}

	return &rd.pageHeader, nil
}

func (rd *Reader) records() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		rd.reader.Reset(rd.pageReader)
		for {
			if err := rd.validateCRC(); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			rc := &record.Record{}
			if err := rc.Read(rd.reader); err != nil {
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

func (rd *Reader) validateCRC() error {
	rb := record.Boundary{}
	if err := rb.Read(rd.reader); err != nil {
		return err
	}

	d, err := rd.reader.Peek(int(rb.Length().Value()))
	if err != nil {
		return err
	}

	if !rb.CRC().Compare(d) {
		return page.ErrInvalidCRC
	}

	return nil
}
