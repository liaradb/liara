package log

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

type LogPageReader struct {
	size       int64
	reader     io.ReadSeeker
	data       []byte
	pageReader *bytes.Reader
	page       LogPageHeader
}

func NewLogPageReader(
	size int64,
	r io.ReadSeeker,
) *LogPageReader {
	body := size - pageHeaderSize
	return &LogPageReader{
		size:   body,
		reader: r,
		data:   make([]byte, body),
	}
}

func (lp *LogPageReader) Seek(pid LogPageID) error {
	_, err := lp.reader.Seek(lp.position(pid, lp.size), io.SeekStart)
	return err
}

// TODO: Should we store this on the header struct?
func (lp *LogPageReader) position(pid LogPageID, size int64) int64 {
	return int64(pid) * (size + pageHeaderSize)
}

func (lpr *LogPageReader) Iterate() iter.Seq2[*LogRecord, error] {
	return lpr.IterateFrom(0)
}

func (lpr *LogPageReader) IterateFrom(pid LogPageID) iter.Seq2[*LogRecord, error] {
	return func(yield func(*LogRecord, error) bool) {
		if err := lpr.Seek(pid); err != nil {
			yield(nil, err)
			return
		}

		for {
			if _, err := lpr.Read(); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			for lr, err := range lpr.Records() {
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(lr, nil) {
					return
				}
			}
		}
	}
}

func (lp *LogPageReader) Read() (*LogPageHeader, error) {
	if err := lp.page.Read(lp.reader); err != nil {
		return nil, err
	}

	// TODO: Do we need to verify read length?
	// TODO: Should we make a new slice?
	if _, err := lp.reader.Read(lp.data); err != nil {
		return nil, err
	}

	lp.initReader()

	return &lp.page, nil
}

func (lp *LogPageReader) initReader() {
	if lp.pageReader == nil {
		lp.pageReader = bytes.NewReader(lp.data)
	} else {
		lp.pageReader.Reset(lp.data)
	}
}

func (lp *LogPageReader) Records() iter.Seq2[*LogRecord, error] {
	r := bufio.NewReader(lp.pageReader)
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
