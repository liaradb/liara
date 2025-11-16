package page

import (
	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/list"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/storage"
)

type BTreePage struct {
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

	BTreePageHeaderSize = NextSize +
		ParentIDSize +
		PrevIDSize +
		NextIDSize +
		LowIDSize
)

func New(data []byte) BTreePage {
	next, data0 := wrap.NewInt32(data)
	parentID, data1 := wrap.NewInt64(data0)
	prevID, data2 := wrap.NewInt64(data1)
	nextID, data3 := wrap.NewInt64(data2)
	lowID, data4 := wrap.NewInt64(data3)

	return BTreePage{
		data:     data,
		next:     next,
		parentID: parentID,
		prevID:   prevID,
		nextID:   nextID,
		lowID:    lowID,
		list:     list.New(data4),
		byteList: bytelist.New(data4),
	}
}

func (p *BTreePage) Append(size int32) (int32, *raw.Buffer, bool) {
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

func (p BTreePage) Length() int32 {
	return int32(len(p.data))
}

func (p BTreePage) Space() int32 {
	return max(p.list.Next()-p.list.Size()-4, 0)
}

func (p BTreePage) hasSpace(size int32) bool {
	s := p.Space()
	return size <= s
}

func (p *BTreePage) ParentID() storage.Offset {
	return storage.Offset(p.parentID.Get())
}

func (p *BTreePage) PrevID() storage.Offset {
	return storage.Offset(p.prevID.Get())
}

func (p *BTreePage) NextID() storage.Offset {
	return storage.Offset(p.nextID.Get())
}

func (p *BTreePage) LowID() storage.Offset {
	return storage.Offset(p.lowID.Get())
}

func (p *BTreePage) SetParentID(o storage.Offset) {
	p.parentID.Set(o.Value())
}

func (p *BTreePage) SetPrevID(o storage.Offset) {
	p.prevID.Set(o.Value())
}

func (p *BTreePage) SetNextID(o storage.Offset) {
	p.nextID.Set(o.Value())
}

func (p *BTreePage) SetLowID(o storage.Offset) {
	p.lowID.Set(o.Value())
}
