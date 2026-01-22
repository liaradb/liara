package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/crclist"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	itemSize = 8
)

type node struct {
	header
	data     []byte
	list     crclist.CRCList
	byteList bytelist.ByteList
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func newNode[S number](size S) *node {
	data := make([]byte, size)
	return newFromSlice(data)
}

func newFromSlice(data []byte) *node {
	header, data0 := newHeader(data)

	return &node{
		header:   header,
		data:     data,
		list:     crclist.New(data0),
		byteList: bytelist.New(data0),
	}
}

// TODO: Test this
func (p *node) Reset(pid action.PageID, tlid action.TimeLineID, rl record.Length) {
	clear(p.data)
	p.list.Clear()
	p.header.Reset(pid, tlid, rl)
}

func (p *node) Append(data []byte) (page.Offset, bool) {
	size := int16(len(data))
	if !p.hasSpace(size) {
		return 0, false
	}

	offset := p.next() - size
	i, ok := p.list.Push(offset, size, page.NewCRC(data))
	if !ok {
		return 0, false
	}

	p.header.setNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return 0, false
	}

	copy(b, data)

	return page.Offset(i), true
}

func (p *node) next() int16 {
	size := p.list.Count()
	if size == 0 {
		return int16(p.list.Length())
	} else {
		return p.header.Next()
	}
}

func (p node) Space() int16 {
	next := p.next()
	size := p.list.Size()
	return max(next-size-itemSize, 0)
}

func (p node) hasSpace(size int16) bool {
	s := p.Space()
	return size <= s
}

func (p node) Child(index int16) ([]byte, bool) {
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

func (p node) Items() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.Items() {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p node) ItemsReverse() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.ItemsReverse() {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p node) ChildrenRange(start, end int16) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.ItemsRange(start, end) {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p *node) Read(r io.ReadSeeker) error {
	_, err := r.Read(p.data)
	if err != nil {
		return err
	}

	p.list.Clear()
	return nil
}

func (p *node) Write(w io.WriteSeeker) error {
	_, err := w.Write(p.data)
	return err
}
