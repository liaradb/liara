package logpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/action"
)

const (
	nextSize = 2

	HeaderSize = 0 +
		page.MagicSize +
		action.TimeLineIDSize
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

func (h *header) reset(
	timeLineID action.TimeLineID,
) {
	h.timeLineID.SetUnsigned(timeLineID.Value())
}

func (h *header) TimeLineID() action.TimeLineID {
	return action.TimeLineID(h.timeLineID.GetUnsigned())
}

func (h header) Size() int {
	return HeaderSize
}

func (h *header) isEmpty() bool {
	return page.Magic(h.magic.Get()).IsEmpty()
}

func (h *header) isPage() bool {
	return page.Magic(h.magic.Get()).IsPage()
}
