package node

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	levelSize  = 1
	highIDSize = 8
	lowIDSize  = 8

	headerSize = 0 +
		page.MagicSize +
		logpage.LogSequenceNumberSize +
		levelSize +
		highIDSize +
		lowIDSize
)

type header struct {
	magic  wrap.Int32
	lsn    wrap.Int64
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)
	lsn, data1 := wrap.NewInt64(data0)
	level, data2 := wrap.NewByte(data1)
	highID, data3 := wrap.NewInt64(data2)
	lowID, data4 := wrap.NewInt64(data3)

	return header{
		magic:  magic,
		lsn:    lsn,
		level:  level,
		highID: highID,
		lowID:  lowID,
	}, data4
}

func (h *header) init() {
	h.magic.Set(int32(page.MagicPage))
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

func (h *header) isEmpty() bool {
	return page.Magic(h.magic.Get()).IsEmpty()
}

func (h *header) isPage() bool {
	return page.Magic(h.magic.Get()).IsPage()
}

func (h *header) LogSequenceNumber() logpage.LogSequenceNumber {
	return logpage.NewLogSequenceNumber(h.lsn.GetUnsigned())
}

func (h *header) SetLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	h.lsn.SetUnsigned(lsn.Value())
}
