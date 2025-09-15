package log

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"iter"

	"github.com/liaradb/liaradb/raw"
)

// # Log Records
//
// ## Common to all
// - prevLSN
// - transID
// - type
//
// ## Update records
// - pageID
// - length
// - offset
// - beforeImage
// - afterImage
//
// # Transaction table
// - pageID
// - recLSN
//
// # Dirty page table
// - transID
// - lastLSN

const (
	BlockSize     = 1024
	SegmentSize   = 1024
	RowHeaderSize = 4 + 4
)

type LogPage struct {
	size            int64
	magic           LogMagic
	id              LogPageID
	timeLineID      TimeLineID
	lengthRemaining LogRecordLength
	data            []byte
	reader          *bytes.Reader
	writer          *bytes.Buffer
	writeBuf        *bufio.Writer
}

func NewLogPage(
	size int64,
) *LogPage {
	writer := bytes.NewBuffer(make([]byte, 0, size))
	return &LogPage{
		size:     size,
		data:     make([]byte, size),
		writer:   writer,
		writeBuf: bufio.NewWriter(writer),
	}
}

func (lp *LogPage) ID() LogPageID                    { return lp.id }
func (lp *LogPage) TimeLineID() TimeLineID           { return lp.timeLineID }
func (lp *LogPage) LengthRemaining() LogRecordLength { return lp.lengthRemaining }
func (lp *LogPage) Data() []byte {
	clear(lp.data)
	copy(lp.data, lp.writer.Bytes())
	return lp.data
}

func (lp *LogPage) Init(id LogPageID, timeLineID TimeLineID) {
	lp.id = id
	lp.timeLineID = timeLineID
}

func (lp *LogPage) reset() {
	lp.writeBuf.Reset(lp.writer)
}

func (lp *LogPage) Append(crc CRC, data []byte) error {
	if err := lp.canInsert(data); err != nil {
		return err
	}

	if err := lp.insert(crc, data); err != nil {
		lp.reset()
		return err
	}

	return nil
}

func (lp *LogPage) canInsert(data []byte) error {
	if lp.recordSize(data) > lp.available() {
		return ErrInsufficientSpace
	}

	return nil
}

func (*LogPage) recordSize(data []byte) int {
	return RowHeaderSize + len(data)
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

	// TODO: Do we need to verify write lengths?
	if n, err := lp.writeBuf.Write(data); err != nil {
		return err
	} else if n != len(data) {
		return raw.ErrOverflow
	}

	return lp.writeBuf.Flush()
}

func (lp *LogPage) Parse(data []byte) error {
	lp.initReader(data)
	return lp.checkMagic()
}

func (lp *LogPage) initReader(data []byte) {
	if lp.reader == nil {
		lp.reader = bytes.NewReader(data)
	} else {
		lp.reader.Reset(data)
	}
}

func (lp *LogPage) checkMagic() error {
	if err := lp.magic.Read(lp.reader); err != nil && err != io.EOF {
		return err
	}

	if lp.magic != LogMagicPage {
		return ErrNotPage
	}

	return nil
}

func (lp *LogPage) Records() iter.Seq2[*LogRecord, error] {
	// lp.initReader()
	return func(yield func(*LogRecord, error) bool) {
		var size uint32
		if err := binary.Read(lp.reader, binary.BigEndian, &size); err != nil {
			yield(nil, err)
			return
		}

		data := make([]byte, size)
		if _, err := lp.reader.Read(data); err != nil {
			yield(nil, err)
			return
		}

		if !yield(&LogRecord{
			data: LogData{data},
		}, nil) {
			return
		}
	}
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

	return nil
}
