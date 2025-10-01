package log

import (
	"bufio"
	"bytes"
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type PageWriter struct {
	bodySize int64
	data     []byte
	writer   *bytes.Buffer
	writeBuf *bufio.Writer
	header   page.PageHeader
}

func newPageWriter(
	size int64,
) *PageWriter {
	body := size - page.PageHeaderSize
	writer := bytes.NewBuffer(make([]byte, 0, body))
	return &PageWriter{
		bodySize: body,
		data:     make([]byte, body),
		writer:   writer,
		writeBuf: bufio.NewWriter(writer),
	}
}

func (pw *PageWriter) ID() page.PageID                    { return pw.header.ID() }
func (pw *PageWriter) TimeLineID() page.TimeLineID        { return pw.header.TimeLineID() }
func (pw *PageWriter) LengthRemaining() page.RecordLength { return pw.header.LengthRemaining() }

// TODO: This is slow
func (pw *PageWriter) Data() []byte {
	// TODO: pw.data is the same backing array as pw.writer
	data := make([]byte, len(pw.data))
	copy(data, pw.writer.Bytes())
	pw.data = data
	return pw.data
}

func (pw *PageWriter) init(id page.PageID, tlid page.TimeLineID, rem page.RecordLength) {
	pw.header = page.NewPageHeader(id, tlid, rem)
}

func (pw *PageWriter) append(crc page.CRC, data []byte) error {
	if !pw.canInsert(data) {
		return ErrInsufficientSpace
	}

	if err := pw.insert(crc, data); err != nil {
		pw.reset()
		return err
	}

	return nil
}

func (pw *PageWriter) reset() {
	pw.writeBuf.Reset(pw.writer)
}

func (pw *PageWriter) canInsert(data []byte) bool {
	return pw.recordSize(data) <= pw.available()
}

func (*PageWriter) recordSize(data []byte) int {
	return page.RecordHeaderSize + len(data)
}

func (pw *PageWriter) available() int {
	return pw.writer.Available()
}

func (pw *PageWriter) insert(crc page.CRC, data []byte) error {
	if err := crc.Write(pw.writeBuf); err != nil {
		return err
	}

	if err := page.NewRecordLength(data).Write(pw.writeBuf); err != nil {
		return err
	}

	if n, err := pw.writeBuf.Write(data); err != nil {
		return err
	} else if n != len(data) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	// TODO: If this fails, should we reset?
	return pw.writeBuf.Flush()
}

func (pw *PageWriter) Flush(w io.WriteSeeker) error {
	if err := pw.seek(w); err != nil {
		return err
	}

	return pw.Write(w)
}

func (pw *PageWriter) seek(w io.Seeker) error {
	_, err := w.Seek(pw.position(), io.SeekStart)
	return err
}

func (pw *PageWriter) position() int64 {
	return pw.header.ID().Size(pw.bodySize + page.PageHeaderSize)
}

func (pw *PageWriter) Write(w io.Writer) error {
	if err := pw.header.Write(w); err != nil {
		return err
	}

	if n, err := w.Write(pw.Data()); err != nil {
		return err
	} else if n < int(pw.bodySize) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	return nil
}

func (pw *PageWriter) SeekTail(r io.Reader) error {
	if err := pw.skipHeader(r); err != nil {
		return err
	}

	if err := pw.loadWriter(r); err != nil {
		return err
	}

	b := bytes.NewBuffer(pw.data)
	for _, err := range pw.records(b) {
		if err != nil {
			return err
		}
	}

	// TODO: Verify this
	l := b.Cap() - (b.Available() + b.Len()) - recordBoundary
	pw.writer = bytes.NewBuffer(pw.data[:l])
	pw.writeBuf = bufio.NewWriter(pw.writer)

	return nil
}

const recordBoundary = 8

func (pw *PageWriter) loadWriter(rd io.Reader) error {
	// TODO: Do we need to verify read length?
	// TODO: Should we handle EOF?
	if _, err := rd.Read(pw.data); err != nil {
		if err != io.EOF {
			return err
		}
	}

	return nil
}

func (pw *PageWriter) skipHeader(rd io.Reader) error {
	data := make([]byte, page.PageHeaderSize)
	// TODO: Do we need to verify read length?
	if _, err := rd.Read(data); err != io.EOF {
		return err
	}

	return nil
}

func (pw *PageWriter) records(rd io.Reader) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {

		for {
			var err error
			if err = pw.skipCRC(rd); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			// TODO: Should we create a new record each time?
			rc := &record.Record{}

			// TODO: Use a buffer
			if err := rc.Read(rd); err != nil {
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

// TODO: We need to rewind the length
func (*PageWriter) skipCRC(r io.Reader) error {
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

	return nil
}
