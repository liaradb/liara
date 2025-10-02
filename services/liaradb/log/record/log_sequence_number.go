package record

import (
	"encoding/binary"
	"io"
)

type LogSequenceNumber uint64

const LogSequenceNumberSize = 8

func (LogSequenceNumber) Size() int { return LogSequenceNumberSize }

func (lsn LogSequenceNumber) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, lsn)
}

func (lsn *LogSequenceNumber) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, lsn)
}
