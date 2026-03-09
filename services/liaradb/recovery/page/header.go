package page

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	nextSize = 2

	headerSize = 0 +
		page.MagicSize +
		nextSize +
		action.PageIDSize +
		action.TimeLineIDSize +
		record.LengthSize
)

type header struct {
	magic           wrap.Int32
	next            wrap.Int16
	id              wrap.Int64
	timeLineID      wrap.Int32
	lengthRemaining wrap.Int32
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)
	next, data1 := wrap.NewInt16(data0)
	id, data2 := wrap.NewInt64(data1)
	tlid, data3 := wrap.NewInt32(data2)
	lr, data4 := wrap.NewInt32(data3)

	return header{
		magic:           magic,
		next:            next,
		id:              id,
		timeLineID:      tlid,
		lengthRemaining: lr,
	}, data4
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

func (h *header) setNext(o int16) {
	h.next.Set(o)
}

func (h *header) ID() action.PageID {
	return action.PageID(h.id.GetUnsigned())
}

func (h *header) LengthRemaining() record.Length {
	return record.NewLength(h.lengthRemaining.GetUnsigned())
}

func (h *header) Next() int16 {
	return h.next.Get()
}

func (h *header) TimeLineID() action.TimeLineID {
	return action.TimeLineID(h.timeLineID.GetUnsigned())
}

func (h header) Size() int {
	return headerSize
}

func (h *header) isEmpty() bool {
	return page.Magic(h.magic.Get()).IsEmpty()
}

func (h *header) isPage() bool {
	return page.Magic(h.magic.Get()).IsPage()
}
