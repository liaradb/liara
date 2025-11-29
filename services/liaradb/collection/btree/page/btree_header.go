package page

import (
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	levelSize  = 1
	highIDSize = 8
	lowIDSize  = 8
	nextSize   = 2

	btreePageHeaderSize = levelSize +
		highIDSize +
		lowIDSize +
		nextSize
)

// TODO: Should we store HighKey?
type btreeHeader struct {
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
	next   wrap.Int16
}

func newHeader(data []byte) (btreeHeader, []byte) {
	level, data0 := wrap.NewByte(data)
	highID, data1 := wrap.NewInt64(data0)
	lowID, data2 := wrap.NewInt64(data1)
	next, data3 := wrap.NewInt16(data2)

	return btreeHeader{
		level:  level,
		highID: highID,
		lowID:  lowID,
		next:   next,
	}, data3
}

func (p *btreeHeader) Level() byte {
	return p.level.GetUnsigned()
}

func (p *btreeHeader) HighID() int64 {
	return p.highID.Get()
}

func (p *btreeHeader) LowID() int64 {
	return p.lowID.Get()
}

func (p *btreeHeader) Next() int16 {
	return p.next.Get()
}

func (p *btreeHeader) setLevel(l byte) {
	p.level.SetUnsigned(l)
}

func (p *btreeHeader) setHighID(o int64) {
	p.highID.Set(o)
}

func (p *btreeHeader) SetLowID(o int64) {
	p.lowID.Set(o)
}

func (p *btreeHeader) setNext(o int16) {
	p.next.Set(o)
}
