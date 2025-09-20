package log

import (
	"bufio"
	"bytes"
	"io"

	"github.com/liaradb/liaradb/log/record"
)

type PageWriter struct {
	size     int64
	data     []byte
	writer   *bytes.Buffer
	writeBuf *bufio.Writer
	header   record.PageHeader
}

func newPageWriter(
	size int64,
) *PageWriter {
	body := size - record.PageHeaderSize
	writer := bytes.NewBuffer(make([]byte, 0, body))
	return &PageWriter{
		size:     body,
		data:     make([]byte, body),
		writer:   writer,
		writeBuf: bufio.NewWriter(writer),
	}
}

func (pw *PageWriter) ID() record.PageID                    { return pw.header.ID() }
func (pw *PageWriter) TimeLineID() record.TimeLineID        { return pw.header.TimeLineID() }
func (pw *PageWriter) LengthRemaining() record.RecordLength { return pw.header.LengthRemaining() }

// TODO: This is slow
func (pw *PageWriter) Data() []byte {
	clear(pw.data)
	copy(pw.data, pw.writer.Bytes())
	return pw.data
}

func (pw *PageWriter) init(id record.PageID, tlid record.TimeLineID, rem record.RecordLength) {
	pw.header = record.NewPageHeader(id, tlid, rem)
}

func (pw *PageWriter) append(crc record.CRC, data []byte) error {
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
	return record.RecordHeaderSize + len(data)
}

func (pw *PageWriter) available() int {
	return int(pw.size) - pw.writer.Len()
}

func (pw *PageWriter) insert(crc record.CRC, data []byte) error {
	if err := crc.Write(pw.writeBuf); err != nil {
		return err
	}

	if err := record.NewRecordLength(data).Write(pw.writeBuf); err != nil {
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
	return pw.header.Position(pw.size)
}

func (pw *PageWriter) Write(w io.Writer) error {
	if err := pw.header.Write(w); err != nil {
		return err
	}

	if n, err := w.Write(pw.Data()); err != nil {
		return err
	} else if n < int(pw.size) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	return nil
}
