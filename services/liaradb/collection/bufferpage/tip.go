package bufferpage

import (
	"context"

	"github.com/liaradb/liaradb/encoder/span"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Tip struct {
	s       *storage.Storage
	fn      link.FileName
	current *BufferPage
	pages   []*BufferPage
	sizes   []int
}

func NewTip(s *storage.Storage, fn link.FileName) Tip {
	return Tip{
		s:  s,
		fn: fn,
	}
}

func (t *Tip) Span(ctx context.Context, size int) (*span.Span, error) {
	s := span.Span{}

	b, err := t.s.RequestCurrent(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	t.current = New(b)

	t.pages = append(t.pages, t.current)

	available := 0
	remaining := size

	p := t.current

	l := t.appendToSpan(&s, p, remaining)

	available = l
	remaining -= l

	for available < size {
		p, err := t.next(ctx)
		if err != nil {
			return nil, err
		}

		l := t.appendToSpan(&s, p, remaining)

		available += l
		remaining -= l
	}

	s.InitIndexes()
	return &s, nil
}

func (t *Tip) appendToSpan(s *span.Span, p *BufferPage, remaining int) int {
	header, data := p.Next(remaining)
	l := len(data)
	t.sizes = append(t.sizes, l)
	if l == 0 {
		return l
	}

	_ = s.Append(header, data)
	return l
}

func (t *Tip) next(ctx context.Context) (*BufferPage, error) {
	b, err := t.s.RequestCurrent(ctx, t.fn)
	if err != nil {
		return nil, err
	}

	p := New(b)
	t.pages = append(t.pages, p)
	return p, nil
}
