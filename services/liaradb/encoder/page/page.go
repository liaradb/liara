package page

import (
	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/slotlist"
)

type Page struct {
	headerSize     int16
	slotHeaderSize int16
	data           []byte
	list           slotlist.SlotList
	byteList       bytelist.ByteList
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
		byteList:       bytelist.New(data[headerSize:]),
	}
}

func (p *Page) Header() []byte {
	return p.data[:p.headerSize]
}

func (p *Page) Slot(i int16) ([]byte, []byte, bool) {
	slot, ok := p.list.Slot(i)
	if !ok {
		return nil, nil, false
	}

	start, end := slot.Range()
	data := p.data[start:end]

	return data[:p.slotHeaderSize], data[p.slotHeaderSize:], true
}

func (p *Page) Next(size int16) ([]byte, []byte) {
	end := p.list.Next()
	start := (end - size) - p.slotHeaderSize
	data := p.data[start:end]

	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) Commit(size int16) (int16, bool) {
	fullSize := size + p.slotHeaderSize
	start := p.list.Next() - fullSize
	i, ok := p.list.Push(start, fullSize)
	if !ok {
		return 0, false
	}

	p.list.SetNext(start)
	return i, true
}
