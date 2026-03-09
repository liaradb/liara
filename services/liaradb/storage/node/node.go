package node

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/crclist"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	itemSize = 8
)

type Node struct {
	header
	buffer   *storage.Buffer
	list     crclist.CRCList
	byteList bytelist.ByteList
}

func New(buffer *storage.Buffer) Node {
	data := buffer.Raw()
	header, data0 := newHeader(data)

	return Node{
		header:   header,
		buffer:   buffer,
		list:     crclist.New(data0),
		byteList: bytelist.New(data0),
	}
}

func (n *Node) Count() int16 { return n.list.Count() }
func (n *Node) Dirty() bool  { return n.buffer.Dirty() }
func (n *Node) Raw() []byte  { return n.buffer.Raw() }

// TODO: Test this
func (n *Node) Release()  { n.buffer.Release() }
func (n *Node) Latch()    { n.buffer.Latch() }
func (n *Node) Unlatch()  { n.buffer.Unlatch() }
func (n *Node) RLatch()   { n.buffer.RLatch() }
func (n *Node) RUnlatch() { n.buffer.RUnlatch() }

func (n *Node) Clear() {
	n.buffer.Clear()
	n.list.Reset()
}

func (n *Node) SetDirty() {
	n.buffer.SetDirty()
}

func (n *Node) Append(size int16, crc page.CRC) (link.RecordPosition, []byte, bool) {
	n.Latch()
	defer n.Unlatch()

	if !n.hasSpace(size) {
		return 0, nil, false
	}

	offset := n.next() - size
	i, ok := n.list.Push(offset, size, crc)
	if !ok {
		return 0, nil, false
	}

	n.header.setNext(offset)
	n.SetDirty()

	b, ok := n.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return link.RecordPosition(i), b, true
}

func (n *Node) Insert(size int16, index int16, crc page.CRC) (link.RecordPosition, []byte, bool) {
	n.Latch()
	defer n.Unlatch()

	if !n.hasSpace(size) {
		return 0, nil, false
	}

	offset := n.next() - size
	i, ok := n.list.Insert(offset, size, crc, index)
	if !ok {
		return 0, nil, false
	}

	n.header.setNext(offset)
	n.SetDirty()

	b, ok := n.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return link.RecordPosition(i), b, true
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

func (n Node) Child(index int16) ([]byte, bool) {
	i, ok := n.list.Item(index)
	if !ok {
		return nil, false
	}

	d, ok := n.byteList.Slice(int64(i.Offset), int64(i.Size))
	if !ok {
		return nil, false
	}

	if !i.CRC.Compare(d) {
		return nil, false
	}

	return d, true
}

// TODO: Should we return old version?
func (n Node) ReplaceChild(index int16, data []byte) bool {
	i, ok := n.list.Item(index)
	if !ok {
		return false
	}

	// Must fit
	if len(data) > int(i.Size) {
		return false
	}

	d, ok := n.byteList.Slice(int64(i.Offset), int64(i.Size))
	if !ok {
		return false
	}

	copy(d, data)
	n.SetDirty()

	return n.list.SetCRC(page.NewCRC(data), index)
}

func (n Node) Children() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range n.list.Items() {
			b, ok := n.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (n Node) ChildrenRange(start, end int16) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range n.list.ItemsRange(start, end) {
			b, ok := n.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}
