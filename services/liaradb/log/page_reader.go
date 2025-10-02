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
	pageSize   int64
	bodySize   int64
	reader     io.ReadSeeker
	data       []byte
	pageReader *bytes.Reader
	pageHeader page.Header
}

func NewPageReader(
	pageSize int64,
	r io.ReadSeeker,
) *PageReader {
	pr := &PageReader{
		pageSize: pageSize,
		reader:   r,
	}
	body := pageSize - int64(pr.pageHeader.Size())
	pr.bodySize = body
	pr.data = make([]byte, body)
	return pr
}

func (pr *PageReader) Seek(pid page.PageID) error {
	_, err := pr.reader.Seek(pid.Size(pr.pageSize), io.SeekStart)
	return err
}

func (pr *PageReader) Iterate() iter.Seq2[*record.Record, error] {
	return pr.IterateFrom(0)
}

// TODO: Test this
func (pr *PageReader) IterateFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		_, err := pr.ReadAt(pid)
		if err != nil {
			yield(nil, err)
			return
		}

		for rc, err := range pr.Records() {
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

func (pr *PageReader) Reverse(size int64) iter.Seq2[*record.Record, error] {
	return pr.ReverseFrom(page.NewActivePageIDFromSize(size, pr.pageSize))
}

// TODO: Change page structure to make reversing easier
// TODO: Test this
func (pr *PageReader) ReverseFrom(pid page.PageID) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		_, err := pr.ReadAt(pid)
		if err != nil {
			yield(nil, err)
			return
		}

		r := list.New()
		for rc, err := range pr.Records() {
			if err != nil {
				yield(nil, err)
				return
			}

			r.PushBack(rc)
		}

		for e := r.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value.(*record.Record), nil) {
				return
			}
		}
	}
}

// TODO: Should we asynchronously prefetch pages?
func (pr *PageReader) Read() (*page.Header, error) {
	if err := pr.pageHeader.Read(pr.reader); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := pr.reader.Read(pr.data); err != nil {
		return nil, err
	}

	pr.initReader()

	return &pr.pageHeader, nil
}

func (pr *PageReader) ReadAt(pid page.PageID) (*page.Header, error) {
	if err := pr.Seek(pid); err != nil {
		return nil, err
	}

	return pr.Read()
}

func (pr *PageReader) initReader() {
	if pr.pageReader == nil {
		pr.pageReader = bytes.NewReader(pr.data)
	} else {
		pr.pageReader.Reset(pr.data)
	}
}

func (pr *PageReader) Records() iter.Seq2[*record.Record, error] {
	r := bufio.NewReader(pr.pageReader)

	return func(yield func(*record.Record, error) bool) {
		for {
			var err error
			// TODO: This reads past the end of the file
			if err = pr.validateCRC(r); err != nil {
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

func (*PageReader) validateCRC(r *bufio.Reader) error {
	var c page.CRC
	if err := c.Read(r); err != nil {
		return err
	}

	rl := page.RecordLength(0)
	if err := rl.Read(r); err != nil {
		return err
	}

	if rl == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(rl))
	if err != nil {
		return err
	}

	if !c.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}
