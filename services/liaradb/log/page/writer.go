package page

import (
	"bufio"
	"bytes"
	"io"
	"iter"

	"github.com/liaradb/liaradb/log/record"
)

type Writer struct {
	bodySize int64
	data     []byte
	writer   *bytes.Buffer
	writeBuf *bufio.Writer
	header   Header
	out      *bytes.Buffer
}

func NewWriter(
	size int64,
) *Writer {
	wr := &Writer{}

	body := size - int64(wr.header.Size())
	writer := bytes.NewBuffer(make([]byte, 0, body))
	wr.bodySize = body
	wr.data = make([]byte, body)
	wr.writer = writer
	wr.writeBuf = bufio.NewWriter(writer)
	wr.out = bytes.NewBuffer(make([]byte, 0, wr.bodySize+int64(wr.header.Size())))

	return wr
}

func (wr *Writer) ID() PageID                     { return wr.header.ID() }
func (wr *Writer) TimeLineID() TimeLineID         { return wr.header.TimeLineID() }
func (wr *Writer) LengthRemaining() record.Length { return wr.header.LengthRemaining() }

// TODO: This is slow
func (wr *Writer) Data() []byte {
	// TODO: pw.data is the same backing array as pw.writer
	data := make([]byte, len(wr.data))
	copy(data, wr.writer.Bytes())
	wr.data = data
	return wr.data
}

func (wr *Writer) Init(id PageID, tlid TimeLineID, rem record.Length) {
	wr.writer.Reset()
	wr.writeBuf.Reset(wr.writer)
	wr.header = NewHeader(id, tlid, rem)
}

func (wr *Writer) Append(rb record.Boundary, data []byte) error {
	if !wr.canInsert(data) {
		return ErrInsufficientSpace
	}

	if err := wr.insert(rb, data); err != nil {
		wr.reset()
		return err
	}

	return nil
}

func (wr *Writer) reset() {
	wr.writeBuf.Reset(wr.writer)
}

func (wr *Writer) canInsert(data []byte) bool {
	return wr.recordSize(data) <= wr.available()
}

func (*Writer) recordSize(data []byte) int {
	return record.BoundarySize + len(data)
}

func (wr *Writer) available() int {
	return wr.writer.Available()
}

func (wr *Writer) insert(rb record.Boundary, data []byte) error {
	if err := rb.Write(wr.writeBuf); err != nil {
		return err
	}

	if n, err := wr.writeBuf.Write(data); err != nil {
		return err
	} else if n != len(data) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	// TODO: If this fails, should we reset?
	return wr.writeBuf.Flush()
}

func (wr *Writer) Flush(w io.WriteSeeker) error {
	if err := wr.seek(w); err != nil {
		return err
	}

	return wr.Write(w)
}

func (wr *Writer) seek(w io.Seeker) error {
	_, err := w.Seek(wr.position(), io.SeekStart)
	return err
}

func (wr *Writer) position() int64 {
	return wr.header.ID().Size(wr.bodySize + int64(wr.header.Size()))
}

func (wr *Writer) Write(w io.Writer) error {
	wr.out.Reset()
	if err := wr.header.Write(wr.out); err != nil {
		return err
	}

	if n, err := wr.out.Write(wr.Data()); err != nil {
		return err
	} else if n < int(wr.bodySize) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	// TODO: Do we need to verify write length?
	_, err := wr.out.WriteTo(w)
	return err
}

func (wr *Writer) SeekTail(r io.Reader) error {
	if err := wr.skipHeader(r); err != nil {
		return err
	}

	if err := wr.loadWriter(r); err != nil {
		return err
	}

	b := bytes.NewBuffer(wr.data)
	for _, err := range wr.records(b) {
		if err != nil {
			return err
		}
	}

	l := wr.bufferSize(b)
	wr.writer = bytes.NewBuffer(wr.data[:l])
	wr.writeBuf = bufio.NewWriter(wr.writer)

	return nil
}

func (wr *Writer) bufferSize(b *bytes.Buffer) int {
	// TODO: Verify this
	return b.Cap() - (b.Available() + b.Len()) - record.BoundarySize
}

func (wr *Writer) loadWriter(rd io.Reader) error {
	// TODO: Do we need to verify read length?
	// TODO: Should we handle EOF?
	if _, err := rd.Read(wr.data); err != nil {
		if err != io.EOF {
			return err
		}
	}

	return nil
}

func (wr *Writer) skipHeader(rd io.Reader) error {
	data := make([]byte, wr.header.Size())
	// TODO: Do we need to verify read length?
	if _, err := rd.Read(data); err != io.EOF {
		return err
	}

	return nil
}

func (wr *Writer) records(rd io.Reader) iter.Seq2[*record.Record, error] {
	rb := &record.Boundary{}
	return func(yield func(*record.Record, error) bool) {

		for {
			var err error
			if err = rb.Read(rd); err != nil {
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
