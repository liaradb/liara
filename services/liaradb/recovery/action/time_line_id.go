package action

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type TimeLineID uint32

const TimeLineIDSize = 4

func (tlid TimeLineID) Value() uint32 {
	return uint32(tlid)
}

func (tlid TimeLineID) Write(w io.Writer) error {
	return raw.WriteInt32(w, tlid)
}

func (tlid *TimeLineID) Read(r io.Reader) error {
	return raw.ReadInt32(r, tlid)
}
