package logpage

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
)

const (
	nextSize = 2

	HeaderSize = 0
)

type LogPage struct {
	page     *page.Page
	handlers []func()
}

func New(size int, slotHeaderSize int) *LogPage {
	page := page.New(size, HeaderSize, slotHeaderSize)
	return &LogPage{
		page: page,
	}
}

func (lp *LogPage) Clear() {
	lp.page.Clear()
}

func (lp *LogPage) Reset() {
	lp.handlers = nil
}

// TODO: How do we use this?
func (lp *LogPage) SetTimeLineID(tlid TimeLineID) {
	lp.page.SetTrackingID(tlid.Value())
}

func (lp *LogPage) TimeLineID() TimeLineID {
	return TimeLineID(lp.page.TrackingID())
}

func (lp *LogPage) Header() []byte {
	return lp.page.Header()
}

func (lp *LogPage) Data() []byte {
	return lp.page.Data()
}

func (lp *LogPage) Replace(r io.Reader) error {
	return lp.page.Replace(r)
}

func (lp *LogPage) Next(size int) ([]byte, []byte) {
	return lp.page.Next(size)
}

func (lp *LogPage) Commit(size int) {
	lp.page.Commit(size)
}

func (lp *LogPage) Fill(data []byte) {
	lp.page.Fill(data)
	lp.handlers = nil
}

func (lp *LogPage) Complete() {
	for _, h := range lp.handlers {
		h()
	}
	lp.handlers = nil
}

func (lp *LogPage) AddHandler(handler func()) {
	if handler != nil {
		lp.handlers = append(lp.handlers, handler)
	}
}

// TODO: Swap receiver and parameter
func (lp *LogPage) Shadow(base *LogPage) {
	lp.page.Fill(base.page.Data())
	lp.handlers = base.handlers
	base.handlers = nil
}

func (lp *LogPage) Slots() iter.Seq2[[]byte, []byte] {
	return lp.page.Slots()
}

func (lp *LogPage) SlotsReverse() iter.Seq2[[]byte, []byte] {
	return lp.page.SlotsReverse()
}
