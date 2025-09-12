package log

import (
	"encoding/binary"
	"io"
)

type LogPageID uint64

func (lpid *LogPageID) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, *lpid)
}

func (lpid *LogPageID) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, lpid)
}
