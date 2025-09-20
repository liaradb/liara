package record

import (
	"encoding/binary"
	"io"
)

type TransactionID uint64

func (tid TransactionID) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, tid)
}

func (tid *TransactionID) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, tid)
}
