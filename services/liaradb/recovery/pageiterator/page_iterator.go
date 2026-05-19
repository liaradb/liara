package pageiterator

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageIterator struct {
	sl *segment.List
}

func New(
	sl *segment.List,
) *PageIterator {
	return &PageIterator{
		sl: sl,
	}
}

func (pi *PageIterator) Forward(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		var p page.Page
		var s record.Span

		for f, err := range pi.sl.IterateFromLSN(lsn) {
			if err != nil {
				yield(nil, err)
				return
			}

			for {
				if err := f.Read(p.Data()); err != nil {
					yield(nil, err)
					return
				}

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

				if !f.NextPage() {
					continue
				}
			}
		}
	}
}
