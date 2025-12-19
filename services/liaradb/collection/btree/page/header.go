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
type header struct {
	level  wrap.Byte
	highID wrap.Int64
	lowID  wrap.Int64
	next   wrap.Int16
}

func newHeader(data []byte) (header, []byte) {
	level, data0 := wrap.NewByte(data)
	highID, data1 := wrap.NewInt64(data0)
	lowID, data2 := wrap.NewInt64(data1)
	next, data3 := wrap.NewInt16(data2)

	return header{
		level:  level,
		highID: highID,
		lowID:  lowID,
		next:   next,
	}, data3
}

func (p *header) Level() byte {
	return p.level.GetUnsigned()
}

func (p *header) HighID() int64 {
	return p.highID.Get()
}

func (p *header) LowID() int64 {
	return p.lowID.Get()
}

func (p *header) Next() int16 {
	return p.next.Get()
}

func (p *header) setLevel(l byte) {
	p.level.SetUnsigned(l)
}

func (p *header) SetHighID(o int64) {
	p.highID.Set(o)
}

func (p *header) SetLowID(o int64) {
	p.lowID.Set(o)
}

func (p *header) setNext(o int16) {
	p.next.Set(o)
}
