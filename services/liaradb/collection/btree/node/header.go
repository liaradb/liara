package node

import (
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	levelSize  = 1
	highIDSize = 8
	lowIDSize  = 8

	headerSize = 0 +
		levelSize +
		highIDSize +
		lowIDSize
)

type header struct {
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
}

func newHeader(data []byte) (header, []byte) {
	level, data0 := wrap.NewByte(data)
	highID, data1 := wrap.NewInt64(data0)
	lowID, data2 := wrap.NewInt64(data1)

	return header{
		level:  level,
		highID: highID,
		lowID:  lowID,
	}, data2
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

func (h *header) setLevel(l byte) {
	h.level.SetUnsigned(l)
}

func (h *header) SetHighID(o link.FilePosition) {
	h.highID.Set(o.Value())
}

func (h *header) SetLowID(o link.FilePosition) {
	h.lowID.Set(o.Value())
}
