package tip

import (
	"context"

	"github.com/liaradb/liaradb/collection/bufferpage"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Tip struct {
	s           *storage.Storage
	fn          link.FileName
	current     *bufferpage.BufferPage
	fileBlockID link.FileBlockID
	slotID      link.SlotID
}

func NewTip(s *storage.Storage, fn link.FileName) Tip {
	return Tip{
		s:  s,
		fn: fn,
	}
}

func (t *Tip) Span(ctx context.Context, size int) (*span.Span, error) {
	s := span.New()

	b, err := t.s.RequestCurrent(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	t.current = bufferpage.New(b, span.FragmentHeaderSize)
	t.fileBlockID = b.FileBlockID()
	t.slotID = link.SlotID(t.current.Count())

	available := 0
	remaining := size

	l := t.appendToSpan(s, t.current, remaining)

	available = l
	remaining -= l

	for available < size {
		p, err := t.next(ctx)
		if err != nil {
			return nil, err
		}

		l := t.appendToSpan(s, p, remaining)

		available += l
		remaining -= l
	}

	s.InitIndexes()
	return s, nil
}

func (t *Tip) appendToSpan(s *span.Span, p *bufferpage.BufferPage, remaining int) int {
	header, data := p.Next(remaining)
	l := len(data)
	if l == 0 {
		p.Release()
		return l
	}

	_ = s.Append(p, 0, header, data)
	return l
}

func (t *Tip) next(ctx context.Context) (*bufferpage.BufferPage, error) {
	b, err := t.s.RequestNext(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	return bufferpage.New(b, span.FragmentHeaderSize), nil
}

func (t *Tip) RecordLocator() link.RecordLocator {
	return link.NewRecordLocator(t.fileBlockID.BlockID(), t.slotID)
}
