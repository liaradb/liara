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
	data     []byte
	list     crclist.CRCList
	byteList bytelist.ByteList
}

func New(buffer *storage.Buffer) Node {
	data := buffer.Raw()
	header, data0 := newHeader(data)

	return Node{
		header:   header,
		buffer:   buffer,
		data:     data,
		list:     crclist.New(data0),
		byteList: bytelist.New(data0),
	}
}

// TODO: Test this
func (p *Node) Clear() {
	p.buffer.Clear()
	p.list.Clear()
}

// TODO: Test this
func (p *Node) Release()  { p.buffer.Release() }
func (p *Node) Latch()    { p.buffer.Latch() }
func (p *Node) Unlatch()  { p.buffer.Unlatch() }
func (p *Node) RLatch()   { p.buffer.RLatch() }
func (p *Node) RUnlatch() { p.buffer.RUnlatch() }

// TODO: Test this
func (p *Node) SetDirty() {
	p.buffer.SetDirty()
}

func (p *Node) Append(size int16, crc page.CRC) (link.RecordPosition, []byte, bool) {
	p.Latch()
	defer p.Unlatch()

	if !p.hasSpace(size) {
		return 0, nil, false
	}

	offset := p.next() - size
	i, ok := p.list.Push(offset, size, crc)
	if !ok {
		return 0, nil, false
	}

	p.header.setNext(offset)
	p.SetDirty()

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return link.RecordPosition(i), b, true
}

func (p *Node) Insert(size int16, index int16, crc page.CRC) (link.RecordPosition, []byte, bool) {
	p.Latch()
	defer p.Unlatch()

	if !p.hasSpace(size) {
		return 0, nil, false
	}

	offset := p.next() - size
	i, ok := p.list.Insert(offset, size, crc, index)
	if !ok {
		return 0, nil, false
	}

	p.header.setNext(offset)
	p.SetDirty()

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return link.RecordPosition(i), b, true
}

func (p *Node) next() int16 {
	size := p.list.Count()
	if size == 0 {
		return int16(p.list.Length())
	} else {
		return p.header.Next()
	}
}

func (p Node) Space() int16 {
	next := p.next()
	size := p.list.Size()
	return max(next-size-itemSize, 0)
}

func (p Node) hasSpace(size int16) bool {
	s := p.Space()
	return size <= s
}

func (p Node) Child(index int16) ([]byte, bool) {
	i, ok := p.list.Item(index)
	if !ok {
		return nil, false
	}

	d, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
	if !ok {
		return nil, false
	}

	if !i.CRC.Compare(d) {
		return nil, false
	}

	return d, true
}

func (p Node) Children() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.Items() {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p Node) ChildrenRange(start, end int16) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.ItemsRange(start, end) {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}
