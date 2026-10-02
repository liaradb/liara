package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/slotlist"
	"github.com/liaradb/liaradb/storage/link"
)

// TODO: Create magic prefix here, instead of in the header.
type Page struct {
	headerSize     int
	slotHeaderSize int
	data           []byte
	body           []byte // TODO: Just use SlotList
	list           slotlist.SlotList
	next           int
}

func New(
	size int,
	headerSize int,
	slotHeaderSize int,
) *Page {
	return NewFromSlice(
		make([]byte, size),
		headerSize,
		slotHeaderSize)
}

func NewFromSlice(
	data []byte,
	headerSize int,
	slotHeaderSize int,
) *Page {
	p := Page{
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		data:           data,
		body:           data[headerSize:],
		list:           slotlist.New(data[headerSize:]),
	}
	p.initNext()
	return &p
}

func (p *Page) Data() []byte {
	return p.data
}

func (p *Page) Length() int16 {
	return int16(len(p.data))
}

func (p *Page) Fill(data []byte) {
	n := copy(p.data, data)
	clear(p.data[n:])
	p.list.Reset()
	p.initNext()
}

func (p *Page) Replace(r io.Reader) error {
	if _, err := r.Read(p.data); err != nil {
		return err
	}

	p.list.Reset()
	p.initNext()
	return nil
}

func (p *Page) Header() []byte {
	return p.data[:p.headerSize]
}

func (p *Page) Slot(i link.SlotID) ([]byte, []byte) {
	slot := p.list.Slot(i)
	h, b := p.slot(slot)
	return h, b
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

func (p *Page) SlotsRange(start, end link.SlotID) iter.Seq2[[]byte, []byte] {
	return func(yield func([]byte, []byte) bool) {
		for slot := range p.list.SlotsRange(start, end) {
			if !yield(p.slot(slot)) {
				return
			}
		}
	}
}

func (p *Page) SlotsReverse() iter.Seq2[[]byte, []byte] {
	return func(yield func([]byte, []byte) bool) {
		for slot := range p.list.SlotsReverse() {
			if !yield(p.slot(slot)) {
				return
			}
		}
	}
}

func (p *Page) slot(s slotlist.Slot) ([]byte, []byte) {
	data := s.SliceUnsafe()
	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) initNext() {
	// TODO: Fix this cast
	p.next = int(p.list.FirstOffset())
}

func (p *Page) Next(size int) ([]byte, []byte) {
	space := p.Space()
	size = min(size, space)
	end := p.next
	start := (end - size) - p.slotHeaderSize
	if end-start <= p.slotHeaderSize {
		return nil, nil
	}

	data := p.body[start:end]
	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) NextMustFit(size int) ([]byte, []byte, bool) {
	space := p.Space()
	if size > space {
		return nil, nil, false
	}

	end := p.next
	start := (end - size) - p.slotHeaderSize
	// TODO: Do we need to test again?
	// if end-start <= p.slotHeaderSize {
	// 	return nil, nil, false
	// }

	data := p.body[start:end]
	return data[:p.slotHeaderSize], data[p.slotHeaderSize:], true
}

func (p *Page) Space() int {
	end := p.next
	// TODO: Fix this cast
	size := int(p.list.NextSize())
	space := (end - size) - p.slotHeaderSize

	return max(space, 0)
}

func (p *Page) Commit(size int) {
	fullSize := size + p.slotHeaderSize
	start := p.next - fullSize

	// TODO: Fix this cast
	_, _ = p.list.Push(int16(start), int16(fullSize))
	p.next = start
}

func (p *Page) Insert(size int, i link.SlotID) {
	fullSize := size + p.slotHeaderSize
	start := p.next - fullSize

	// TODO: Fix this cast
	_, _ = p.list.Insert(int16(start), int16(fullSize), i)
	p.next = start
}

func (p *Page) Clear() {
	clear(p.data)
	p.list.Clear()
	p.initNext()
}

func (p *Page) Count() link.SlotID {
	return p.list.Count()
}
