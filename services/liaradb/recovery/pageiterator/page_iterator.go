package pageiterator

import (
	"container/list"
	"io"
	"iter"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/span"
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

// This is unused.  It may be useful for in-place Recover.
func (pi *PageIterator[R]) Forward(lsn logpage.LogSequenceNumber) iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		var s span.Span
		l := list.New()

		for sg, err := range pi.sl.IterateFromLSN(lsn) {
			if err != nil {
				var r R
				yield(r, err)
				return
			}

			size := sg.Size()
			if size == 0 {
				return
			}

			for p, err := range pi.iterateSegment(sg, size) {
				if err != nil {
					var r R
					yield(r, err)
					return
				}

				e := l.PushBack(p)

				freePages := false
				freeAll := false
				for h, d := range p.Slots() {
					f := s.Append(h, d)
					if freeAll {
						freeAll = false
					}

					if f.Index() == 0 {
						if rc, err := pi.ToRecord(s); !yield(rc, err) || err != nil {
							pi.pool.Put(p)
							return
						}

						freePages = true
						freeAll = true
						s = span.Span{}
					}
				}

				if freePages {
					for n := e.Prev(); n != nil; n = n.Prev() {
						p := n.Value.(*logpage.LogPage)
						l.Remove(n)
						pi.pool.Put(p)
					}
				}
				if freeAll {
					p := e.Value.(*logpage.LogPage)
					l.Remove(e)
					pi.pool.Put(p)
				}
			}
		}
	}
}

func (pi *PageIterator[R]) Reverse() iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		var s span.Span

		for sg, err := range pi.sl.Reverse() {
			if err != nil {
				var r R
				yield(r, err)
				return
			}

			if size, err := sg.SeekTail(); err != nil {
				var r R
				yield(r, err)
				return
			} else if size == 0 {
				continue
			}

			// TODO: Return Page
			for p, err := range pi.reverseSegment(sg) {
				if err != nil {
					if err != io.EOF {
						var r R
						yield(r, err)
					}
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
			}
		}
	}
}

func (pi *PageIterator[R]) iterateSegment(sg *segment.File, size int64) iter.Seq2[*logpage.LogPage, error] {
	return func(yield func(*logpage.LogPage, error) bool) {
		for {
			if p, err := pi.getPage(sg); !yield(p, err) ||
				err != nil ||
				!sg.NextPageUntilSize(size) {
				return
			}
		}
	}
}

func (pi *PageIterator[R]) reverseSegment(sg *segment.File) iter.Seq2[*logpage.LogPage, error] {
	return func(yield func(*logpage.LogPage, error) bool) {
		for {
			if p, err := pi.getPage(sg); !yield(p, err) ||
				err != nil ||
				!sg.PrevPageUntilStart() {
				return
			}
		}
	}
}

func (pi *PageIterator[R]) getPage(sg *segment.File) (*logpage.LogPage, error) {
	ok := false
	p := pi.pool.Get()
	defer func() {
		if !ok {
			pi.pool.Put(p)
		}
	}()

	if err := p.Replace(sg); err != nil {
		return nil, err
	}

	ok = true
	return p, nil
}

func (pi *PageIterator[R]) ToRecord(s span.Span) (R, error) {
	rc := pi.create()
	if err := rc.Read(s); err != nil {
		var r R
		return r, err
	}

	return rc, nil
}
