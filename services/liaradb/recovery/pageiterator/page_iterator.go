package pageiterator

import (
	"iter"

	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageIterator struct {
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

			for {
				p := pi.pool.Get()
				if err := f.Read(p.Data()); err != nil {
					yield(nil, err)
					return
				}

				p.Reset()
				for h, d := range p.Slots() {
					f := record.NewFragment(h, d)
					s.Append(f)

					if f.Index() == 0 {
						rc := &record.Record{}
						if err := rc.Read(s); err != nil {
							yield(nil, err)
							return
						}

						if !yield(rc, nil) {
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
