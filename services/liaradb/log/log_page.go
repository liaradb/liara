package log

import (
	"bytes"
	"encoding/binary"
	"io"
	"iter"
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
	BlockSize   uint64 = 1024
	SegmentSize uint64 = 1024
)

type LogPage struct {
	magic           LogMagic
	id              LogPageID
	timeLineID      TimeLineID
	lengthRemaining LogRecordLength
	positions       []uint32
	data            []byte
	reader          *bytes.Reader
}

func NewLogPage(
	size int64,
	id LogPageID,
	timeLineID TimeLineID,
) *LogPage {
	return &LogPage{
		id:         id,
		timeLineID: timeLineID,
		data:       make([]byte, size),
	}
}

func (lp *LogPage) ID() LogPageID                    { return lp.id }
func (lp *LogPage) TimeLineID() TimeLineID           { return lp.timeLineID }
func (lp *LogPage) LengthRemaining() LogRecordLength { return lp.lengthRemaining }
func (lp *LogPage) Data() []byte                     { return lp.data }

func (lp *LogPage) Parse(data []byte) error {
	lp.data = data
	lp.initReader()
	return lp.checkMagic()
}

func (lp *LogPage) initReader() {
	if lp.reader == nil {
		lp.reader = bytes.NewReader(lp.data)
	} else {
		lp.reader.Reset(lp.data)
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
	lp.initReader()
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
