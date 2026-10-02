package node

import (
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	levelSize  = 1
	highIDSize = 8
	lowIDSize  = 8

	headerSize = 0 +
		logpage.LogSequenceNumberSize +
		levelSize +
		highIDSize +
		lowIDSize
)

type header struct {
	lsn    wrap.Int64
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
}

func newHeader(data []byte) (header, []byte) {
	lsn, data0 := wrap.NewInt64(data)
	level, data1 := wrap.NewByte(data0)
	highID, data2 := wrap.NewInt64(data1)
	lowID, data3 := wrap.NewInt64(data2)

	return header{
		lsn:    lsn,
		level:  level,
		highID: highID,
		lowID:  lowID,
	}, data3
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

func (h *header) LogSequenceNumber() logpage.LogSequenceNumber {
	return logpage.NewLogSequenceNumber(h.lsn.GetUnsigned())
}

func (h *header) setLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	h.lsn.SetUnsigned(lsn.Value())
}
