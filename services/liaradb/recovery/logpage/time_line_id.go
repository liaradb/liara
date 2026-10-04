package logpage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type TimeLineID uint64

const TimeLineIDSize = 8

func (tlid TimeLineID) Value() uint64 {
	return uint64(tlid)
}

func (tlid TimeLineID) Write(w io.Writer) error {
	return raw.WriteInt64(w, tlid)
}

func (tlid *TimeLineID) Read(r io.Reader) error {
	return raw.ReadInt64(r, tlid)
}
