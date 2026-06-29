package logpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type LogPage struct {
	*page.Page
	header  header
	handler func()
}

func (lp *LogPage) Clear(lsn record.LogSequenceNumber) {
	lp.Page.Clear()
	lp.header.init()
}

func (lp *LogPage) Reset() {
	lp.handler = nil
}

func New(size int16, slotHeaderSize int16) *LogPage {
	page := page.New(size, HeaderSize, slotHeaderSize)
	header, _ := newHeader(page.Header())
	return &LogPage{
		Page:   page,
		header: header,
	}
}

func (lp *LogPage) Fill(data []byte, h func()) {
	lp.Page.Fill(data)
	lp.SetHandler(h)
}

func (lp *LogPage) Complete() {
	if lp.handler != nil {
		lp.handler()
		lp.handler = nil
	}
}

func (lp *LogPage) Handler() func() {
	return lp.handler
}

func (lp *LogPage) SetHandler(handler func()) {
	lp.handler = handler
}
