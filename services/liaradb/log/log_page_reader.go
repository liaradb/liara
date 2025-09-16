package log

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

type LogPageReader struct {
	size   int64
	data   []byte
	reader *bytes.Reader
	page   LogPageHeader
}

func newLogPageReader(
	size int64,
) *LogPageReader {
	body := size - PageHeaderSize
	return &LogPageReader{
		size: body,
		data: make([]byte, body),
	}
}

func (lp *LogPageReader) Seek(w io.WriteSeeker, pid LogPageID) error {
	_, err := w.Seek(lp.position(pid, lp.size), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (lp *LogPageReader) position(pid LogPageID, size int64) int64 {
	return int64(pid) * (size + PageHeaderSize)
}

func (lp *LogPageReader) Read(r io.Reader) (*LogPageHeader, error) {
	if err := lp.page.Read(r); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := r.Read(lp.data); err != nil {
		return nil, err
	}

	lp.initReader()

	return &lp.page, nil
}

func (lp *LogPageReader) initReader() {
	if lp.reader == nil {
		lp.reader = bytes.NewReader(lp.data)
	} else {
		lp.reader.Reset(lp.data)
	}
}

func (lp *LogPageReader) Records() iter.Seq2[*LogRecord, error] {
	r := bufio.NewReader(lp.reader)
	lr := &LogRecord{}

	return func(yield func(*LogRecord, error) bool) {
		for {
			var err error
			if err = lp.validateCRC(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			// TODO: Use a buffer
			if err := lr.Read(r); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(lr, nil) {
				return
			}
		}
	}
}

func (*LogPageReader) validateCRC(r *bufio.Reader) error {
	var c CRC
	if err := c.Read(r); err != nil {
		return err
	}

	lrl := LogRecordLength(0)
	if err := lrl.Read(r); err != nil {
		return err
	}

	if lrl == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(lrl))
	if err != nil {
		return err
	}

	if !c.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}
