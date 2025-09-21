package page

import (
	"encoding/binary"
	"io"
)

type RecordLength uint32

const RecordLengthSize = 4

func NewRecordLength(d []byte) RecordLength {
	return RecordLength(len(d))
}

func (rl RecordLength) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, rl)
}

func (rl *RecordLength) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, rl)
}
