package page

import (
	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/tuplelist"
)

const (
	itemSize = 4
)

type header = btreeHeader

type BTreePage struct {
	header
	data     []byte
	list     tuplelist.TupleList
	byteList bytelist.ByteList
}

func New(data []byte) BTreePage {
	header, data0 := newHeader(data)

	return BTreePage{
		header:   header,
		data:     data,
		list:     tuplelist.New(data0),
		byteList: bytelist.New(data0),
	}
}

func (p *BTreePage) Append(size int16) (int16, *raw.Buffer, bool) {
	if !p.hasSpace(size) {
		return 0, nil, false
	}

	offset := p.next() - size
	i, ok := p.list.Push(offset, size) // TODO: Is this the correct order?
	if !ok {
		return 0, nil, false
	}

	p.header.setNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size)) // TODO: Is this the correct order?
	if !ok {
		return 0, nil, false
	}

	return i, b, true
}

func (p BTreePage) Length() int16 {
	return int16(len(p.data))
}

func (p *BTreePage) next() int16 {
	size := p.list.Count()
	if size == 0 {
		return int16(p.list.Length())
	} else {
		return p.header.Next()
	}
}

func (p BTreePage) Space() int16 {
	next := p.next()
	size := p.list.Size()
	return max(next-size-itemSize, 0)
}

func (p BTreePage) hasSpace(size int16) bool {
	s := p.Space()
	return size <= s
}
