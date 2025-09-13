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
	Magic           LogMagic
	ID              LogPageID
	TimeLineID      TimeLineID
	LengthRemaining int
	Positions       []uint32
	Data            []byte
	Reader          *bytes.Reader
}

func NewLogPage(size int64) *LogPage {
	return &LogPage{
		Data: make([]byte, size),
	}
}

func (lp *LogPage) Parse(data []byte) error {
	lp.Data = data
	lp.initReader()
	return lp.checkMagic()
}

func (lp *LogPage) initReader() {
	if lp.Reader == nil {
		lp.Reader = bytes.NewReader(lp.Data)
	} else {
		lp.Reader.Reset(lp.Data)
	}
}

func (lp *LogPage) checkMagic() error {
	if err := lp.Magic.Read(lp.Reader); err != nil && err != io.EOF {
		return err
	}

	if lp.Magic != LogMagicPage {
		return ErrNotPage
	}

	return nil
}

func (lp *LogPage) Records() iter.Seq2[*LogRecord, error] {
	lp.initReader()
	return func(yield func(*LogRecord, error) bool) {
		var size uint32
		if err := binary.Read(lp.Reader, binary.BigEndian, &size); err != nil {
			yield(nil, err)
			return
		}

		data := make([]byte, size)
		if _, err := lp.Reader.Read(data); err != nil {
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
