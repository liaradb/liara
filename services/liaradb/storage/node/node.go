package node

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/tuplelist"
	"github.com/liaradb/liaradb/storage"
)

const (
	itemSize = 4
)

type Node struct {
	header
	buffer   *storage.Buffer
	data     []byte
	list     tuplelist.TupleList
	byteList bytelist.ByteList
}

func New(buffer *storage.Buffer) Node {
	data := buffer.Raw()
	header, data0 := newHeader(data)

	return Node{
		header:   header,
		buffer:   buffer,
		data:     data,
		list:     tuplelist.New(data0),
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

func (p *Node) Append(size int16) (int16, []byte, bool) {
	if !p.hasSpace(size) {
		return 0, nil, false
	}

	offset := p.next() - size
	i, ok := p.list.Push(offset, size)
	if !ok {
		return 0, nil, false
	}

	p.header.setNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return i, b, true
}

func (p *Node) Insert(size int16, index int16) (int16, []byte, bool) {
	if !p.hasSpace(size) {
		return 0, nil, false
	}

	offset := p.next() - size
	i, ok := p.list.Insert(offset, size, index)
	if !ok {
		return 0, nil, false
	}

	p.header.setNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, nil, false
	}

	return i, b, true
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
	offset, size, ok := p.list.Item(index)
	if !ok {
		return nil, false
	}

	return p.byteList.Slice(int64(offset), int64(size))
}

func (p Node) Children() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for offset, size := range p.list.Items() {
			b, ok := p.byteList.Slice(int64(offset), int64(size))
			if !ok || !yield(b) {
				return
			}
		}
	}
}

func (p Node) ChildrenRange(start, end int16) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for offset, size := range p.list.ItemsRange(start, end) {
			b, ok := p.byteList.Slice(int64(offset), int64(size))
			if !ok || !yield(b) {
				return
			}
		}
	}
}
