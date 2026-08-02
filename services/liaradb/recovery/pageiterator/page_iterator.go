package pageiterator

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageIterator[R Record] struct {
	sl     *segment.List
	pool   *pagepool.PagePool
	create func() R
}

type Record interface {
	Read(r io.Reader) error
}

func New[R Record](
	sl *segment.List,
	pl *pagepool.PagePool,
	create func() R,
) *PageIterator[R] {
	return &PageIterator[R]{
		sl:     sl,
		pool:   pl,
		create: create,
	}
}

// TODO: This is not used.  It may be useful for in-place Recover.
func (pi *PageIterator[R]) Forward(lsn logpage.LogSequenceNumber) iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		var s span.Span

		for f, err := range pi.sl.IterateFromLSN(lsn) {
			if err != nil {
				var r R
				yield(r, err)
				return
			}

			size := f.Size()
			if size == 0 {
				return
			}

			for {
				// TODO: Return Page
				p := pi.pool.Get()
				if err := p.Replace(f); err != nil {
					var r R
					yield(r, err)
					pi.pool.Put(p)
					return
				}

				for h, d := range p.Slots() {
					f := s.Append(h, d)

					if f.Index() == 0 {
						if rc, err := pi.ToRecord(s); !yield(rc, err) || err != nil {
							pi.pool.Put(p)
							return
						}

						s = span.Span{}
					}
				}

				if !f.NextPageUntilSize(size) {
					break
				}
			}
		}
	}
}

func (pi *PageIterator[R]) Reverse() iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		var s span.Span

		for f, err := range pi.sl.Reverse() {
			if err != nil {
				var r R
				yield(r, err)
				return
			}

			if size, err := f.SeekTail(); err != nil {
				var r R
				yield(r, err)
				return
			} else if size == 0 {
				continue
			}

			for {
				// TODO: Return Page
				p := pi.pool.Get()
				if err := p.Replace(f); err != nil {
					if err != io.EOF {
						var r R
						yield(r, err)
					}
					pi.pool.Put(p)
					return
				}

				for h, d := range p.SlotsReverse() {
					f := s.Append(h, d)

					if f.Count()-1 == f.Index() {
						s.Reverse()

						if rc, err := pi.ToRecord(s); !yield(rc, err) || err != nil {
							pi.pool.Put(p)
							return
						}

						s = span.Span{}
					}
				}

				if !f.PrevPageUntilStart() {
					break
				}
			}
		}
	}
}

func (pi *PageIterator[R]) ToRecord(s span.Span) (R, error) {
	rc := pi.create()
	if err := rc.Read(s); err != nil {
		var r R
		return r, err
	}

	return rc, nil
}
