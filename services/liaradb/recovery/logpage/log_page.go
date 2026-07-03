package logpage

import (
	"slices"

	"github.com/liaradb/liaradb/encoder/page"
)

type LogPage struct {
	*page.Page
	header
	handlers []func()
}

func (lp *LogPage) Clear() {
	lp.Page.Clear()
	lp.header.init()
}

func (lp *LogPage) Reset() {
	lp.handlers = nil
}

func New(size int, slotHeaderSize int) *LogPage {
	page := page.New(size, HeaderSize, slotHeaderSize)
	header, _ := newHeader(page.Header())
	return &LogPage{
		Page:   page,
		header: header,
	}
}

func (lp *LogPage) Fill(data []byte) {
	lp.Page.Fill(data)
	lp.handlers = nil
}

func (lp *LogPage) Complete() {
	for _, h := range lp.handlers {
		h()
	}
	lp.handlers = nil
}

func (lp *LogPage) Handlers() []func() {
	return lp.handlers
}

func (lp *LogPage) AddHandler(handler func()) {
	if handler != nil {
		lp.handlers = append(lp.handlers, handler)
	}
}

func (lp *LogPage) Copy(base *LogPage) {
	lp.Page.Fill(base.Data())
	lp.handlers = slices.Clone(base.handlers)
}
