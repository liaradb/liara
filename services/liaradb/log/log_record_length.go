package log

import (
	"encoding/binary"
	"io"
)

type LogRecordLength uint32

func NewLogRecordLength(d []byte) LogRecordLength {
	return LogRecordLength(len(d))
}

func (lrl LogRecordLength) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, lrl)
}

func (lrl *LogRecordLength) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, lrl)
}
