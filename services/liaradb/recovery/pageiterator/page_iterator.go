package pageiterator

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageIterator struct {
	size int16
	sl   *segment.List
	pool pagequeue.Pool
}

func New(
	sl *segment.List,
	size int16,
	headerSize int16,
	slotHeaderSize int16,
) *PageIterator {
	return &PageIterator{
		size: size,
		sl:   sl,
		pool: pagequeue.NewPool(size, headerSize, slotHeaderSize),
	}
}

func (pi *PageIterator) Forward(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		var s record.Span

		for f, err := range pi.sl.IterateFromLSN(lsn) {
			if err != nil {
				yield(nil, err)
				return
			}

			size, err := f.Size()
			if err != nil {
				yield(nil, err)
				return
			}

			if size == 0 {
				return
			}

			for {
				// TODO: Return Page
				p := pi.pool.Get()
				if err := p.Replace(f); err != nil {
					yield(nil, err)
					return
				}

				for h, d := range p.Slots() {
					f := record.NewFragment(h, d)
					s.Append(f)

					if f.Index() == 0 {
						if rc, err := s.ToRecord(); !yield(rc, err) || err != nil {
							return
						}

						s = record.NewSpan()
					}
				}

				if !f.NextPageUntilSize(size) {
					break
				}
			}
		}
	}
}

func (pi *PageIterator) Reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		var s record.Span

		for f, err := range pi.sl.Reverse() {
			if err != nil {
				yield(nil, err)
				return
			}

			size, err := f.Size()
			if err != nil {
				yield(nil, err)
				return
			}

			if size == 0 {
				continue
			}

			// TODO: Combine these
			if err := f.SeekTail(); err != nil {
				yield(nil, err)
				return
			}

			for {
				// TODO: Return Page
				p := pi.pool.Get()
				if err := p.Replace(f); err != nil {
					if err != io.EOF {
						yield(nil, err)
					}
					return
				}
				p.SlotsReverse()

				for h, d := range p.SlotsReverse() {
					f := record.NewFragment(h, d)
					s.Append(f)

					if f.Count()-1 == f.Index() {
						s.Reverse()

						if rc, err := s.ToRecord(); !yield(rc, err) || err != nil {
							return
						}

						s = record.NewSpan()
					}
				}

				if !f.PrevPageUntilStart() {
					break
				}
			}
		}
	}
}

func (pi *PageIterator) position(pid action.PageID) int64 {
	return pid.Position(int64(pi.size))
}
