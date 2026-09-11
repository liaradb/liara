package node

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/slotlist"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/slice"
)

const (
	itemSize = 4
)

type Node struct {
	header
	buffer *storage.Buffer
	data   []byte
	body   []byte
	list   slotlist.SlotList
}

func New(buffer *storage.Buffer) Node {
	data := buffer.Raw()
	header, data0 := newHeader(data)

	if header.isEmpty() {
		header.init()
	}

	return Node{
		header: header,
		buffer: buffer,
		data:   data,
		body:   data0,
		list:   slotlist.New(data0),
	}
}

func (n *Node) Clear() {
	n.buffer.Clear()
	n.list.Clear()
	n.header.init()
}

func (n *Node) IsPage() bool { return n.header.isPage() }

// TODO: Test this
func (n *Node) Release()  { n.buffer.Release() }
func (n *Node) Latch()    { n.buffer.Latch() }
func (n *Node) Unlatch()  { n.buffer.Unlatch() }
func (n *Node) RLatch()   { n.buffer.RLatch() }
func (n *Node) RUnlatch() { n.buffer.RUnlatch() }

func (n *Node) SetDirty() {
	n.buffer.SetDirty()
}

func (n *Node) SetLevel(l byte) {
	n.header.setLevel(l)
}

func (n *Node) Append(size int16) (int16, []byte, bool) {
	if !n.hasSpace(size) {
		return 0, nil, false
	}

	offset := n.next() - size
	i, ok := n.list.Push(offset, size)
	if !ok {
		return 0, nil, false
	}

	n.header.setNext(offset)

	b, ok := n.slice(offset, size)
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return i, b, true
}

// TODO: Change to SlotID
func (n *Node) Insert(size int16, index int16) (int16, []byte, bool) {
	if !n.hasSpace(size) {
		return 0, nil, false
	}

	offset := n.next() - size
	i, ok := n.list.Insert(offset, size, link.SlotID(index))
	if !ok {
		return 0, nil, false
	}

	n.header.setNext(offset)

	b, ok := n.slice(offset, size)
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return i, b, true
}

func (n Node) Length() int16 {
	return int16(len(n.data))
}

func (n Node) Count() int16 {
	return n.list.Count()
}

func (n *Node) next() int16 {
	size := n.list.Count()
	if size == 0 {
		return int16(n.list.Length())
	} else {
		return n.header.Next()
	}
}

func (n Node) Space() int16 {
	next := n.next()
	size := n.list.Size()
	return max(next-size-itemSize, 0)
}

func (n Node) hasSpace(size int16) bool {
	s := n.Space()
	return size <= s
}

// TODO: Change to SlotID
func (n Node) Child(index int16) ([]byte, bool) {
	slot, ok := n.list.Slot(link.SlotID(index))
	if !ok {
		return nil, false
	}

	return n.slice(slot.Offset(), slot.Size())
}

func (n Node) Children() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for slot := range n.list.Slots() {
			b, ok := n.slice(slot.Offset(), slot.Size())
			if !ok || !yield(b) {
				return
			}
		}
	}
}

// TODO: Change to SlotID
func (n Node) ChildrenRange(start, end int16) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for slot := range n.list.SlotsRange(link.SlotID(start), link.SlotID(end)) {
			b, ok := n.slice(slot.Offset(), slot.Size())
			if !ok || !yield(b) {
				return
			}
		}
	}
}

func (n Node) slice(offset, size int16) ([]byte, bool) {
	return slice.Slice(n.body, int64(offset), int64(size))
}
