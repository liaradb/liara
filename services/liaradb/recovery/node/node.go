package node

import (
	"errors"
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/crclist"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/action"
	p "github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	itemSize = 8
)

type Node struct {
	header
	data     []byte
	list     crclist.CRCList
	byteList bytelist.ByteList
}

var _ p.Page = (*Node)(nil)

func New(data []byte) Node {
	header, data0 := newHeader(data)

	return Node{
		header:   header,
		data:     data,
		list:     crclist.New(data0),
		byteList: bytelist.New(data0),
	}
}

// TODO: Test this
func (p *Node) Clear() {
	clear(p.data)
	p.list.Clear()
}

func (p *Node) Reset(pid action.PageID, tlid action.TimeLineID, rl record.Length) {
	p.header.Reset(pid, tlid, rl)
}

func (p *Node) Add(data []byte) (page.Offset, error) {
	size := int16(len(data))
	if !p.hasSpace(size) {
		return 0, errors.New("no space")
	}

	offset := p.next() - size
	i, ok := p.list.Push(offset, size, page.NewCRC(data))
	if !ok {
		return 0, errors.New("no space")
	}

	p.header.setNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, errors.New("no space")
	}

	copy(b, data)

	return page.Offset(i), nil
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

func (p Node) Items() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.Items() {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p Node) ItemsReverse() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.ItemsReverse() {
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

func (p *Node) Read(r io.ReadSeeker) error {
	_, err := r.Read(p.data)
	return err
}

func (p *Node) Write(w io.WriteSeeker) error {
	_, err := w.Write(p.data)
	return err
}
