package logpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type LogPage struct {
	*page.Page
	lsn     record.LogSequenceNumber
	handler func()
}

func (lp *LogPage) Clear(lsn record.LogSequenceNumber) {
	lp.Page.Clear()
	lp.lsn = lsn
}

func (lp *LogPage) Reset() {
	lp.handler = nil
}

func NewLogPage(size int16, headerSize int16, slotHeaderSize int16) *LogPage {
	return &LogPage{
		Page: page.New(size, headerSize, slotHeaderSize),
	}
}

func (lp *LogPage) Fill(data []byte, h func()) {
	lp.Page.Fill(data)
	lp.SetHandler(h)
}

func (lp *LogPage) LSN() record.LogSequenceNumber { return lp.lsn }

func (lp *LogPage) Complete() {
	if lp.handler != nil {
		lp.handler()
		lp.handler = nil
	}
}

func (lp *LogPage) SetLSN(lsn record.LogSequenceNumber) {
	lp.lsn = lsn
}

func (lp *LogPage) Handler() func() {
	return lp.handler
}

func (lp *LogPage) SetHandler(handler func()) {
	lp.handler = handler
}
