package page

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type TimeLineID uint32

const timeLineIDSize = 4

func (tlid TimeLineID) Write(w io.Writer) error {
	return raw.WriteInt32(w, tlid)
}

func (tlid *TimeLineID) Read(r io.Reader) error {
	return raw.ReadInt32(r, tlid)
}
