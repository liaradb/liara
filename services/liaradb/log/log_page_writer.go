package log

import (
	"bufio"
	"bytes"
	"io"
)

const recordHeaderSize = crcSize + recordLengthSize

type LogPageWriter struct {
	size     int64
	data     []byte
	writer   *bytes.Buffer
	writeBuf *bufio.Writer
	header   LogPageHeader
}

func newLogPageWriter(
	size int64,
) *LogPageWriter {
	body := size - pageHeaderSize
	writer := bytes.NewBuffer(make([]byte, 0, body))
	return &LogPageWriter{
		size:     body,
		data:     make([]byte, body),
		writer:   writer,
		writeBuf: bufio.NewWriter(writer),
	}
}

func (lp *LogPageWriter) ID() LogPageID                 { return lp.header.ID() }
func (lp *LogPageWriter) TimeLineID() TimeLineID        { return lp.header.TimeLineID() }
func (lp *LogPageWriter) LengthRemaining() RecordLength { return lp.header.LengthRemaining() }

// TODO: This is slow
func (lp *LogPageWriter) Data() []byte {
	clear(lp.data)
	copy(lp.data, lp.writer.Bytes())
	return lp.data
}

func (lp *LogPageWriter) init(id LogPageID, tlid TimeLineID, rem RecordLength) {
	lp.header = newLogPageHeader(id, tlid, rem)
}

func (lp *LogPageWriter) append(crc CRC, data []byte) error {
	if !lp.canInsert(data) {
		return ErrInsufficientSpace
	}

	if err := lp.insert(crc, data); err != nil {
		lp.reset()
		return err
	}

	return nil
}

func (lp *LogPageWriter) reset() {
	lp.writeBuf.Reset(lp.writer)
}

func (lp *LogPageWriter) canInsert(data []byte) bool {
	return lp.recordSize(data) <= lp.available()
}

func (*LogPageWriter) recordSize(data []byte) int {
	return recordHeaderSize + len(data)
}

func (lp *LogPageWriter) available() int {
	return int(lp.size) - lp.writer.Len()
}

func (lp *LogPageWriter) insert(crc CRC, data []byte) error {
	if err := crc.Write(lp.writeBuf); err != nil {
		return err
	}

	if err := NewRecordLength(data).Write(lp.writeBuf); err != nil {
		return err
	}

	if n, err := lp.writeBuf.Write(data); err != nil {
		return err
	} else if n != len(data) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	// TODO: If this fails, should we reset?
	return lp.writeBuf.Flush()
}

func (lp *LogPageWriter) Flush(w io.WriteSeeker) error {
	if err := lp.seek(w); err != nil {
		return err
	}

	return lp.Write(w)
}

func (lp *LogPageWriter) seek(w io.Seeker) error {
	_, err := w.Seek(lp.position(), io.SeekStart)
	return err
}

func (lp *LogPageWriter) position() int64 {
	return lp.header.position(lp.size)
}

func (lp *LogPageWriter) Write(w io.Writer) error {
	if err := lp.header.Write(w); err != nil {
		return err
	}

	if n, err := w.Write(lp.Data()); err != nil {
		return err
	} else if n < int(lp.size) {
		// TODO: Do we need to verify write length?
		return io.ErrShortWrite
	}

	return nil
}
