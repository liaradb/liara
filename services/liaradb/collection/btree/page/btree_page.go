package page

import (
	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/list"
	"github.com/liaradb/liaradb/encoder/raw"
)

const (
	itemSize = 2
)

type header = btreeHeader

type BTreePage struct {
	header
	data     []byte
	list     list.List
	byteList bytelist.ByteList
}

func New(data []byte) BTreePage {
	header, data0 := newHeader(data)

	return BTreePage{
		header:   header,
		data:     data,
		list:     list.New(data0),
		byteList: bytelist.New(data0),
	}
}

func (p *BTreePage) Append(size int16) (int16, *raw.Buffer, bool) {
	if !p.hasSpace(size) {
		return 0, nil, false
	}

	offset := p.list.Next() - size
	i, ok := p.list.Push(offset)
	if !ok {
		return 0, nil, false
	}

	p.list.SetNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok {
		return 0, nil, false
	}

	return i, b, true
}

func (p BTreePage) Length() int16 {
	return int16(len(p.data))
}

func (p BTreePage) Space() int16 {
	return max(p.list.Next()-p.list.Size()-itemSize, 0)
}

func (p BTreePage) hasSpace(size int16) bool {
	s := p.Space()
	return size <= s
}
