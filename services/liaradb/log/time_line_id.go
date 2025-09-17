package log

import (
	"encoding/binary"
	"io"
)

type TimeLineID uint32

const timeLineIDSize = 4

func (tlid TimeLineID) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, tlid)
}

func (tlid *TimeLineID) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, tlid)
}
