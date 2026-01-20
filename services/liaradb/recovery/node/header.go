package node

import (
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

// TODO: Should this have a magic entry?
const (
	nextSize = 2

	headerSize = 0 +
		nextSize +
		action.PageIDSize +
		action.TimeLineIDSize +
		record.LengthSize
)

type header struct {
	next            wrap.Int16
	id              wrap.Int64
	timeLineID      wrap.Int32
	lengthRemaining wrap.Int32
}

func newHeader(data []byte) (header, []byte) {
	next, data0 := wrap.NewInt16(data)
	id, data1 := wrap.NewInt64(data0)
	tlid, data2 := wrap.NewInt32(data1)
	lr, data3 := wrap.NewInt32(data2)

	return header{
		next:            next,
		id:              id,
		timeLineID:      tlid,
		lengthRemaining: lr,
	}, data3
}

func (h *header) Reset(
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
