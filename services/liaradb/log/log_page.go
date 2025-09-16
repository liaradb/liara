package log

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

type LogPage struct {
	size            int64
	id              LogPageID
	timeLineID      TimeLineID
	lengthRemaining LogRecordLength
	data            []byte
	reader          *bytes.Reader
	writer          *bytes.Buffer
	writeBuf        *bufio.Writer
}

func newLogPage(
	size int64,
) *LogPage {
	body := size - PageHeaderSize
	writer := bytes.NewBuffer(make([]byte, 0, body))
	return &LogPage{
		size:     body,
		data:     make([]byte, body),
		writer:   writer,
		writeBuf: bufio.NewWriter(writer),
	}
}

func (lp *LogPage) ID() LogPageID                    { return lp.id }
func (lp *LogPage) TimeLineID() TimeLineID           { return lp.timeLineID }
func (lp *LogPage) LengthRemaining() LogRecordLength { return lp.lengthRemaining }

// TODO: This is slow
func (lp *LogPage) Data() []byte {
	clear(lp.data)
	copy(lp.data, lp.writer.Bytes())
	return lp.data
}

func (lp *LogPage) init(id LogPageID, timeLineID TimeLineID) {
	lp.id = id
	lp.timeLineID = timeLineID
}

func (lp *LogPage) append(crc CRC, data []byte) error {
	if !lp.canInsert(data) {
		return ErrInsufficientSpace
	}

	if err := lp.insert(crc, data); err != nil {
		lp.reset()
		return err
	}

	return nil
}

func (lp *LogPage) reset() {
	lp.writeBuf.Reset(lp.writer)
}

func (lp *LogPage) canInsert(data []byte) bool {
	return lp.recordSize(data) <= lp.available()
}

func (*LogPage) recordSize(data []byte) int {
	return RecordHeaderSize + len(data)
}

func (lp *LogPage) available() int {
	return int(lp.size) - lp.writer.Len()
}

func (lp *LogPage) insert(crc CRC, data []byte) error {
	if err := crc.Write(lp.writeBuf); err != nil {
		return err
	}

	if err := NewLogRecordLength(data).Write(lp.writeBuf); err != nil {
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

func (lp *LogPage) Flush(w interface {
	io.Writer
	io.Seeker
}) error {
	if err := lp.Seek(w); err != nil {
		return err
	}

	return lp.Write(w)
}

func (lp *LogPage) Seek(w io.WriteSeeker) error {
	_, err := w.Seek(lp.position(), io.SeekStart)
	return err
}

func (lp *LogPage) position() int64 {
	return int64(lp.id) * (lp.size + PageHeaderSize)
}

func (lp *LogPage) Write(w io.Writer) error {
	if err := LogMagicPage.Write(w); err != nil {
		return err
	}

	if err := lp.id.Write(w); err != nil {
		return err
	}

	if err := lp.timeLineID.Write(w); err != nil {
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

func (lp *LogPage) Read(r io.Reader) error {
	if err := LogMagicPage.ReadIsPage(r); err != nil {
		return err
	}

	if err := lp.id.Read(r); err != nil {
		return err
	}

	if err := lp.timeLineID.Read(r); err != nil {
		return err
	}

	// TODO: Do we need to verify read length?
	if _, err := r.Read(lp.data); err != nil {
		return err
	}

	lp.initReader()

	return nil
}

func (lp *LogPage) initReader() {
	if lp.reader == nil {
		lp.reader = bytes.NewReader(lp.data)
	} else {
		lp.reader.Reset(lp.data)
	}
}

func (lp *LogPage) Records() iter.Seq2[*LogRecord, error] {
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

func (*LogPage) validateCRC(r *bufio.Reader) error {
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
