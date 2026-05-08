package pagequeue

import (
	"github.com/liaradb/liaradb/encoder/slotlist"
)

type Page struct {
	headerSize     int16
	slotHeaderSize int16
	data           []byte
	list           slotlist.SlotList
}

func NewPage(
	size int,
	headerSize int16,
	slotHeaderSize int16,
) *Page {
	data := make([]byte, size)
	return &Page{
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		data:           data,
		list:           slotlist.New(data[headerSize:]),
	}
}

func (p *Page) Header() []byte {
	return p.data[:p.headerSize]
}

func (p *Page) Slot(i int16) ([]byte, []byte, bool) {
	item, ok := p.list.Item(i)
	if !ok {
		return nil, nil, false
	}

	start := p.headerSize + item.Offset
	end := start + item.Size
	slot := p.data[start:end]

	return slot[p.slotHeaderSize:], slot[:p.slotHeaderSize], true
}
