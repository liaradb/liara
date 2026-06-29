package logpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	nextSize = 2

	HeaderSize = 0 +
		page.MagicSize +
		action.PageIDSize +
		action.TimeLineIDSize +
		record.LengthSize
)

type header struct {
	magic           wrap.Int32
	id              wrap.Int64
	timeLineID      wrap.Int32
	lengthRemaining wrap.Int32
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)
	id, data1 := wrap.NewInt64(data0)
	tlid, data2 := wrap.NewInt32(data1)
	lr, data3 := wrap.NewInt32(data2)

	return header{
		magic:           magic,
		id:              id,
		timeLineID:      tlid,
		lengthRemaining: lr,
	}, data3
}

func (h *header) init() {
	h.magic.Set(int32(page.MagicPage))
}

func (h *header) reset(
	id action.PageID,
	timeLineID action.TimeLineID,
	lengthRemaining record.Length,
) {
	h.id.SetUnsigned(id.Value())
	h.timeLineID.SetUnsigned(timeLineID.Value())
	h.lengthRemaining.SetUnsigned(lengthRemaining.Value())
}

func (h *header) ID() action.PageID {
	return action.PageID(h.id.GetUnsigned())
}

func (h *header) LengthRemaining() record.Length {
	return record.NewLength(h.lengthRemaining.GetUnsigned())
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
