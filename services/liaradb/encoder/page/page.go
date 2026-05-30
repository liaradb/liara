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
	next           int16
}

func New(
	size int16,
	headerSize int16,
	slotHeaderSize int16,
) *Page {
	return NewFromSlice(
		make([]byte, size),
		headerSize,
		slotHeaderSize)
}

func NewFromSlice(
	data []byte,
	headerSize int16,
	slotHeaderSize int16,
) *Page {
	p := Page{
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		data:           data,
		list:           slotlist.New(data[headerSize:]),
		byteList:       bytelist.New(data[headerSize:]),
	}
	p.initNext()
	return &p
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

	h, b := p.slot(slot)
	return h, b, true
}

func (p *Page) Slots() iter.Seq2[[]byte, []byte] {
	return func(yield func([]byte, []byte) bool) {
		for slot := range p.list.Slots() {
			if !yield(p.slot(slot)) {
				return
			}
		}
	}
}

func (p *Page) slot(s slotlist.Slot) ([]byte, []byte) {
	start, end := s.Range()
	data := p.data[start:end]

	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) initNext() {
	if last, ok := p.list.Last(); ok {
		p.next = last.Offset()
	} else {
		// TODO: Fix this cast
		p.next = int16(p.list.Length())
	}
}

func (p *Page) Next(size int16) ([]byte, []byte) {
	space := p.space()
	size = min(size, space)
	end := p.next
	start := (end - size) - p.slotHeaderSize
	data := p.data[start:end]
	if end-start < p.slotHeaderSize {
		return nil, nil
	}

	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) space() int16 {
	end := p.next
	size := p.list.Size()
	start := (end - size) - p.slotHeaderSize

	return max(start, 0)
}

func (p *Page) Commit(size int16) (int16, bool) {
	fullSize := size + p.slotHeaderSize
	start := p.next - fullSize
	i, ok := p.list.Push(start, fullSize)
	if !ok {
		return 0, false
	}

	p.next = start
	return i, true
}

func (p *Page) Clear() {
	clear(p.data)
	p.list.Clear()
	p.initNext()
}
