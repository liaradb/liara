package node

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	itemSize = 4
)

type Node struct {
	header
	page   *page.Page
	buffer *storage.Buffer
}

func New(buffer *storage.Buffer) Node {
	page := page.NewFromSlice(buffer.Raw(), headerSize, 0)
	header, _ := newHeader(page.Header())

	if header.isEmpty() {
		header.init()
	}

	return Node{
		header: header,
		page:   page,
		buffer: buffer,
	}
}

func (n *Node) Clear() {
	// n.buffer.Clear()
	n.page.Clear()
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

func (n *Node) Append(size int16) ([]byte, bool) {
	// TODO: Fix this cast
	_, b, ok := n.page.NextMustFit(int(size))
	return b, ok
}

func (n *Node) Commit(size int16) bool {
	// TODO: Fix this cast
	return n.page.Commit(int(size))
}

func (n *Node) Insert(size int16, index link.SlotID) ([]byte, bool) {
	// TODO: Fix this cast
	_, b, ok := n.page.NextMustFit(int(size))
	if !ok {
		return nil, false
	}

	if !n.page.Insert(int(size), index) {
		return nil, false
	}

	return b, ok
}

func (n *Node) Length() int16 {
	return n.page.Length()
}

func (n *Node) Count() link.SlotID {
	return n.page.Count()
}

func (n *Node) Space() int16 {
	return int16(n.page.Space())
}

func (n *Node) Child(index link.SlotID) ([]byte, bool) {
	_, b, ok := n.page.Slot(index)
	return b, ok
}

func (n *Node) Children() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for _, b := range n.page.Slots() {
			if !yield(b) {
				return
			}
		}
	}
}

func (n *Node) ChildrenRange(start, end link.SlotID) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for _, b := range n.page.SlotsRange(start, end) {
			if !yield(b) {
				return
			}
		}
	}
}
