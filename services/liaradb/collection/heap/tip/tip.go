package tip

import (
	"context"
	"math"

	"github.com/liaradb/liaradb/collection/heap/bufferpage"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/bufferlist"
)

type Tip struct {
	l  log.Log
	fn link.FileName
}

func NewTip(l log.Log, fn link.FileName) Tip {
	return Tip{
		l:  l,
		fn: fn,
	}
}

func (t *Tip) Record(ctx context.Context, size int) (*record.Record, error) {
	r := record.New()
	ok := false
	defer func() {
		if !ok {
			r.Release()
		}
	}()

	b, err := t.l.PinCurrent(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	available := 0
	remaining := size

	l := t.appendToRecord(r, b, t.capRemaining(remaining))

	available = l
	remaining -= l

	for available < size {
		b, err := t.l.PinNext(ctx, t.fn)
		if err != nil {
			return nil, err
		}

		l := t.appendToRecord(r, b, t.capRemaining(remaining))

		available += l
		remaining -= l
	}

	r.InitIndexes()
	ok = true
	return r, nil
}

func (*Tip) capRemaining(remaining int) int16 {
	return int16(min(remaining, math.MaxInt16))
}

func (t *Tip) appendToRecord(r *record.Record, b *bufferlist.Buffer, remaining int16) int {
	ok := false
	p := bufferpage.New(b, t.l)
	defer func() {
		if !ok {
			p.Release()
		}
	}()

	header, data := p.Next(remaining)
	l := len(data)
	if l == 0 {
		return l
	}

	_ = r.Append(p, link.SlotID(p.Count()), header, data)
	ok = true
	return l
}
