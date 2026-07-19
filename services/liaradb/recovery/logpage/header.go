package logpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	nextSize = 2

	HeaderSize = 0 +
		page.MagicSize +
		TimeLineIDSize
)

type header struct {
	magic      wrap.Int32
	timeLineID wrap.Int32
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)
	tlid, data1 := wrap.NewInt32(data0)

	return header{
		magic:      magic,
		timeLineID: tlid,
	}, data1
}

func (h *header) init() {
	h.magic.Set(int32(page.MagicPage))
}

// TODO: How do we use this?
func (h *header) SetTimeLineID(tlid TimeLineID) {
	h.timeLineID.SetUnsigned(tlid.Value())
}

func (h *header) TimeLineID() TimeLineID {
	return TimeLineID(h.timeLineID.GetUnsigned())
}

func (h *header) IsEmpty() bool {
	return page.Magic(h.magic.Get()).IsEmpty()
}

func (h *header) IsPage() bool {
	return page.Magic(h.magic.Get()).IsPage()
}
