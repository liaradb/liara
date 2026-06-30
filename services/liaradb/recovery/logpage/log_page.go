package logpage

import "github.com/liaradb/liaradb/encoder/page"

type LogPage struct {
	*page.Page
	header  header
	handler func()
}

func (lp *LogPage) Clear() {
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

func (lp *LogPage) Fill(data []byte) {
	lp.Page.Fill(data)
	lp.SetHandler(nil)
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

func (lp *LogPage) Copy(base *LogPage) {
	lp.Page.Fill(base.Data())
	lp.handler = base.handler
}
