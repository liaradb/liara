package node

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	levelSize  = 1
	highIDSize = 8
	lowIDSize  = 8
	nextSize   = 2

	headerSize = 0 +
		page.MagicSize +
		levelSize +
		highIDSize +
		lowIDSize +
		nextSize
)

// TODO: Should we store HighKey?
type header struct {
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
	next   wrap.Int16
}

func newHeader(data []byte) (header, []byte) {
	// TODO: Should this have a magic entry?
	_, data0 := wrap.NewInt32(data) // Magic
	level, data1 := wrap.NewByte(data0)
	highID, data2 := wrap.NewInt64(data1)
	lowID, data3 := wrap.NewInt64(data2)
	next, data4 := wrap.NewInt16(data3)

	return header{
		level:  level,
		highID: highID,
		lowID:  lowID,
		next:   next,
	}, data4
}

func (h *header) Level() byte {
	return h.level.GetUnsigned()
}

func (h *header) HighID() link.FilePosition {
	return link.FilePosition(h.highID.Get())
}

func (h *header) LowID() link.FilePosition {
	return link.FilePosition(h.lowID.Get())
}

func (h *header) Next() int16 {
	return h.next.Get()
}

func (h *header) setLevel(l byte) {
	h.level.SetUnsigned(l)
}

func (h *header) SetHighID(o link.FilePosition) {
	h.highID.Set(o.Value())
}

func (h *header) SetLowID(o link.FilePosition) {
	h.lowID.Set(o.Value())
}

func (h *header) setNext(o int16) {
	h.next.Set(o)
}
