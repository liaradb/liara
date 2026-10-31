package node

import (
	"slices"

	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	levelSize  = 1
	highIDSize = 8
	lowIDSize  = 8

	HeaderSize = 0 +
		levelSize +
		highIDSize +
		lowIDSize
)

type header struct {
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
	base   []byte
	data   []byte
	dirty  bool
}

func newHeader(data []byte) (header, []byte) {
	level, data0 := wrap.NewByte(data)
	highID, data1 := wrap.NewInt64(data0)
	lowID, data2 := wrap.NewInt64(data1)

	return header{
		level:  level,
		highID: highID,
		lowID:  lowID,
		base:   slices.Clone(data),
		data:   data,
	}, data2
}

func (h *header) Level() byte {
	return h.level.GetUnsigned()
}

func (h *header) HighID() link.BlockID {
	return link.BlockID(h.highID.Get())
}

func (h *header) LowID() link.BlockID {
	return link.BlockID(h.lowID.Get())
}

func (h *header) setLevel(l byte) {
	h.level.SetUnsigned(l)
	h.dirty = true
}

func (h *header) SetHighID(o link.BlockID) {
	h.highID.Set(o.Value())
	h.dirty = true
}

func (h *header) SetLowID(o link.BlockID) {
	h.lowID.Set(o.Value())
	h.dirty = true
}

func (h *header) commit() {
	copy(h.base, h.data)
	h.dirty = false
}
