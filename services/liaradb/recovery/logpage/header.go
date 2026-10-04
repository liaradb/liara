package logpage

import (
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	nextSize = 2

	HeaderSize = 0 +
		TimeLineIDSize
)

type header struct {
	timeLineID wrap.Int64
}

func newHeader(data []byte) (header, []byte) {
	tlid, data0 := wrap.NewInt64(data)

	return header{
		timeLineID: tlid,
	}, data0
}

// TODO: How do we use this?
func (h *header) SetTimeLineID(tlid TimeLineID) {
	h.timeLineID.SetUnsigned(tlid.Value())
}

func (h *header) TimeLineID() TimeLineID {
	return TimeLineID(h.timeLineID.GetUnsigned())
}
