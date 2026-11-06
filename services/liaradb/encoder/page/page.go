package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/slotlist"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	headerSize = 0 +
		MagicSize +
		trackingIDSize
	trackingIDSize = 8
)

type Page struct {
	magic          wrap.Int32
	trackingID     wrap.Int64
	headerSize     int
	slotHeaderSize int16
	data           []byte
	body           []byte // TODO: Just use SlotList
	list           slotlist.SlotList
	next           int16
}

func New(
	size int,
	headerSize int,
	slotHeaderSize int16,
) *Page {
	return NewFromSlice(
		make([]byte, size),
		headerSize,
		slotHeaderSize)
}

func NewFromSlice(
	data []byte,
	headerSize int,
	slotHeaderSize int16,
) *Page {
	magic, data0 := wrap.NewInt32(data)
	trackingID, data1 := wrap.NewInt64(data0)
	p := Page{
		magic:          magic,
		trackingID:     trackingID,
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		data:           data,
		body:           data1[headerSize:],
		list:           slotlist.New(data1[headerSize:]),
	}
	p.initMagic() // TODO: Should this be called here?
	p.initNext()
	return &p
}

func (p *Page) Data() []byte {
	return p.data
}

func (p *Page) Fill(data []byte) {
	n := copy(p.data, data)
	p.reset(n)
}

func (p *Page) Replace(r io.Reader) error {
	n, err := r.Read(p.data)
	if err != nil {
		return err
	}

	p.reset(n)
	return nil
}

func (p *Page) reset(n int) {
	// TODO: If n < len(p.data) should we panic?
	clear(p.data[n:])
	p.list.Reset()
	p.initNext()
}

func (p *Page) Header() []byte {
	return p.data[headerSize : headerSize+p.headerSize]
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
	p.next = p.list.FirstOffset()
}

func (p *Page) Next(size int16) ([]byte, []byte) {
	data, ok := p.nextSlot(size)
	if !ok {
		return nil, nil
	}

	return data[:p.slotHeaderSize], data[p.slotHeaderSize:]
}

func (p *Page) nextSlot(size int16) ([]byte, bool) {
	space := p.Space()
	size = min(size, space)
	end := p.next
	start := (end - size) - p.slotHeaderSize
	if end-start <= p.slotHeaderSize {
		return nil, false
	}

	return p.body[start:end], true
}

func (p *Page) recoverSlot(size int16) []byte {
	space := p.Space()
	size = min(size, space)
	end := p.next
	start := (end - size)
	return p.body[start:end]
}

func (p *Page) NextMustFit(size int16) ([]byte, []byte, bool) {
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

func (p *Page) Space() int16 {
	end := p.next
	size := p.list.NextSize()
	space := (end - size) - p.slotHeaderSize

	return max(space, 0)
}

func (p *Page) CommitAppend(size int16) link.SlotID {
	fullSize := size + p.slotHeaderSize
	start := p.next - fullSize

	_, sid := p.list.Push(start, fullSize)
	p.next = start
	return sid
}

func (p *Page) Insert(size int16, i link.SlotID) link.SlotID {
	fullSize := size + p.slotHeaderSize
	start := p.next - fullSize

	_, sid := p.list.Insert(start, fullSize, i)
	p.next = start
	return sid
}

func (p *Page) RecoverAppend(data []byte) {
	slot := p.recoverSlot(int16(len(data)))
	if n := copy(slot, data); n < len(data) {
		panic("incomplete recovery")
	}

	p.CommitAppend(int16(len(data)) - p.slotHeaderSize)
}

func (p *Page) RevertAppend() {
	clear(p.list.Pop().SliceUnsafe())
	p.initNext()
}

func (p *Page) RecoverUpdate(i link.SlotID, data []byte) {
	if n := copy(p.list.Slot(i).SliceUnsafe(), data); n < len(data) {
		panic("incomplete recovery")
	}
}

func (p *Page) RecoverInsert(i link.SlotID, data []byte) {
	slot := p.recoverSlot(int16(len(data)))
	if n := copy(slot, data); n < len(data) {
		panic("incomplete recovery")
	}

	_ = p.Insert(int16(len(data))-p.slotHeaderSize, i)
}

func (p *Page) RecoverHeader(data []byte) {
	if n := copy(p.Header(), data); n < len(data) {
		panic("incomplete recovery")
	}
}

func (p *Page) RecoverClear() {
	p.Clear()
}

func (p *Page) Clear() {
	clear(p.data)
	p.list.Clear()
	p.initMagic()
	p.initNext()
}

func (p *Page) Count() link.SlotID {
	return p.list.Count()
}

// TODO: Make sure to call this
func (p *Page) initMagic() {
	p.magic.Set(int32(MagicPage))
}

// TODO: We are not using this
func (p *Page) isEmpty() bool {
	return Magic(p.magic.Get()).IsEmpty()
}

// TODO: We are not using this
func (p *Page) isPage() bool {
	return Magic(p.magic.Get()).IsPage()
}

func (h *Page) TrackingID() uint64              { return h.trackingID.GetUnsigned() }
func (h *Page) SetTrackingID(trackingID uint64) { h.trackingID.SetUnsigned(trackingID) }
