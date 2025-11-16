package page

import (
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/storage"
)

const (
	levelSize    = 1
	parentIDSize = 8
	prevIDSize   = 8
	nextIDSize   = 8
	lowIDSize    = 8
	nextSize     = 4

	btreePageHeaderSize = levelSize +
		parentIDSize +
		prevIDSize +
		nextIDSize +
		lowIDSize +
		nextSize
)

type btreeHeader struct {
	level    wrap.Byte
	parentID wrap.Int64
	prevID   wrap.Int64
	nextID   wrap.Int64
	lowID    wrap.Int64
	next     wrap.Int32
}

func newHeader(data []byte) (btreeHeader, []byte) {
	level, data0 := wrap.NewByte(data)
	parentID, data1 := wrap.NewInt64(data0)
	prevID, data2 := wrap.NewInt64(data1)
	nextID, data3 := wrap.NewInt64(data2)
	lowID, data4 := wrap.NewInt64(data3)
	next, data5 := wrap.NewInt32(data4)

	return btreeHeader{
		level:    level,
		parentID: parentID,
		prevID:   prevID,
		nextID:   nextID,
		lowID:    lowID,
		next:     next,
	}, data5
}

func (p *btreeHeader) Level() byte {
	return p.level.GetUnsigned()
}

func (p *btreeHeader) ParentID() storage.Offset {
	return storage.Offset(p.parentID.Get())
}

func (p *btreeHeader) PrevID() storage.Offset {
	return storage.Offset(p.prevID.Get())
}

func (p *btreeHeader) NextID() storage.Offset {
	return storage.Offset(p.nextID.Get())
}

func (p *btreeHeader) LowID() storage.Offset {
	return storage.Offset(p.lowID.Get())
}

func (p *btreeHeader) setLevel(l byte) {
	p.level.SetUnsigned(l)
}

func (p *btreeHeader) setParentID(o storage.Offset) {
	p.parentID.Set(o.Value())
}

func (p *btreeHeader) setPrevID(o storage.Offset) {
	p.prevID.Set(o.Value())
}

func (p *btreeHeader) setNextID(o storage.Offset) {
	p.nextID.Set(o.Value())
}

func (p *btreeHeader) setLowID(o storage.Offset) {
	p.lowID.Set(o.Value())
}
