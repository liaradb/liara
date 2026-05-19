package pageiterator

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageIterator struct {
	list *segment.List
}

func New(
	list *segment.List,
) *PageIterator {
	return &PageIterator{
		list: list,
	}
}

func (pi *PageIterator) Forward(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		var p page.Page
		var s record.Span
		var rc *record.Record = &record.Record{}
		for f, err := range pi.list.IterateFromLSN(lsn) {
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
				}

				if err := rc.Read(s); err != nil {
					yield(nil, err)
					return
				}

				if !yield(rc, nil) {
					return
				}

				if !f.NextPage() {
					continue
				}
			}
		}
	}
}
