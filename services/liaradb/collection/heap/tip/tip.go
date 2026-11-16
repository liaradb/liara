package tip

import (
	"context"
	"math"

	"github.com/liaradb/liaradb/collection/heap/bufferpage"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Tip struct {
	s  *storage.Storage
	fn link.FileName
}

func NewTip(s *storage.Storage, fn link.FileName) Tip {
	return Tip{
		s:  s,
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

	b, err := t.s.RequestCurrent(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	available := 0
	remaining := size

	p := bufferpage.New(b)
	l := t.appendToRecord(r, p, t.capRemaining(remaining))

	available = l
	remaining -= l

	for available < size {
		p, err := t.next(ctx)
		if err != nil {
			return nil, err
		}

		l := t.appendToRecord(r, p, t.capRemaining(remaining))

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

func (t *Tip) appendToRecord(r *record.Record, p *bufferpage.BufferPage, remaining int16) int {
	ok := false
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

func (t *Tip) next(ctx context.Context) (*bufferpage.BufferPage, error) {
	b, err := t.s.RequestNext(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	return bufferpage.New(b), nil
}
