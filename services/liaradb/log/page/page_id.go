package page

import (
	"encoding/binary"
	"io"
)

type PageID uint64

const pageIDSize = 8

func (pid PageID) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, pid)
}

func (pid *PageID) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, pid)
}
