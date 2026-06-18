package pageiterator

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/span"
)

type PageIterator struct {
	size int16
	sl   *segment.List
	pool pagepool.PagePool
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
		pool: pagepool.New(size, headerSize, slotHeaderSize),
	}
}

func (pi *PageIterator) Forward(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		var s span.Span

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
					f := s.Append(h, d)

					if f.Index() == 0 {
						if rc, err := pi.ToRecord(s); !yield(rc, err) || err != nil {
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

func (pi *PageIterator) Reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		var s span.Span

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
			if _, err := f.SeekTail(); err != nil {
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
					f := s.Append(h, d)

					if f.Count()-1 == f.Index() {
						s.Reverse()

						if rc, err := pi.ToRecord(s); !yield(rc, err) || err != nil {
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

func (*PageIterator) ToRecord(s span.Span) (*record.Record, error) {
	rc := record.Record{}
	if err := rc.Read(s); err != nil {
		return nil, err
	}

	return &rc, nil
}
