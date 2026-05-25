package page

import (
	"iter"

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

func New(
	size int16,
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

func NewFromSlice(
	data []byte,
	headerSize int16,
	slotHeaderSize int16,
) *Page {
	return &Page{
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		data:           data,
		list:           slotlist.New(data[headerSize:]),
		byteList:       bytelist.New(data[headerSize:]),
	}
}

func (p *Page) Data() []byte {
	return p.data
}

func (p *Page) Fill(data []byte) {
	n := copy(p.data, data)
	clear(p.data[n:])
}

func (p *Page) Reset() {
	p.list.Reset()
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

func (p *Page) Slots() iter.Seq2[[]byte, []byte] {
	return func(yield func([]byte, []byte) bool) {
		for slot := range p.list.Slots() {
			start, end := slot.Range()
			data := p.data[start:end]

			if !yield(data[:p.slotHeaderSize], data[p.slotHeaderSize:]) {
				return
			}
		}
	}
}

func (p *Page) Next(size int16) ([]byte, []byte) {
	space := p.space()
	size = min(size, space)
	end := p.list.Next()
	start := (end - size) - p.slotHeaderSize
	data := p.data[start:end]
	if end-start < p.slotHeaderSize {
		return nil, nil
	}

	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) space() int16 {
	end := p.list.Next()
	size := p.list.Size()
	start := (end - size) - p.slotHeaderSize

	return max(start, 0)
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

func (p *Page) Clear() {
	clear(p.data)
	p.list.Clear()
}
