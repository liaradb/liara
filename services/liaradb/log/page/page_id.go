package page

import (
	"encoding/binary"
	"io"
)

type PageID uint64

const pageIDSize = 8

func NewPageIDFromSize(size int64, pageSize int64) PageID {
	if pageSize == 0 {
		return 0
	}
	pid := size / pageSize
	if size%pageSize != 0 {
		pid++
	}
	return PageID(pid)
}

func (pid PageID) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, pid)
}

func (pid *PageID) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, pid)
}
