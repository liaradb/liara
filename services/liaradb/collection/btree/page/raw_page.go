package page

import (
	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/list"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/wrap"
)

type RawPage struct {
	data     []byte
	next     wrap.Int32
	parentID wrap.Int64
	prevID   wrap.Int64
	nextID   wrap.Int64
	lowID    wrap.Int64
	list     list.List
	byteList bytelist.ByteList
}

const (
	NextSize     = 4
	ParentIDSize = 8
	PrevIDSize   = 8
	NextIDSize   = 8
	LowIDSize    = 8

	RawPageHeaderSize = NextSize +
		ParentIDSize +
		PrevIDSize +
		NextIDSize +
		LowIDSize
)

func New(data []byte) RawPage {
	next, data0 := wrap.NewInt32(data)
	_, data1 := wrap.NewInt64(data0)
	_, data2 := wrap.NewInt64(data1)
	_, data3 := wrap.NewInt64(data2)
	_, data4 := wrap.NewInt64(data3)

	return RawPage{
		data:     data,
		next:     next,
		list:     list.New(data4),
		byteList: bytelist.New(data4),
	}
}

func (p *RawPage) Append(size int32) (int32, *raw.Buffer, bool) {
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

func (p RawPage) Length() int32 {
	return int32(len(p.data))
}

func (p RawPage) Space() int32 {
	return max(p.list.Next()-p.list.Size()-4, 0)
}

func (p RawPage) hasSpace(size int32) bool {
	s := p.Space()
	return size <= s
}
