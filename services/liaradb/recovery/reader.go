package recovery

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/segmentio"
	"github.com/liaradb/liaradb/util/iterator"
)

type reader struct {
	sl *segment.List
	sr *segmentio.Reader
	it *pageiterator.PageIterator
}

func newReader(
	pageSize int64,
	sl *segment.List,
	sl2 *segment.List,
) *reader {
	return &reader{
		sl: sl,
		sr: segmentio.NewReader(pageSize),
		it: pageiterator.New(sl2, int16(pageSize), page.HeaderSize, page.ItemHeaderSize),
	}
}

// Iterate in reverse until record type.
//
// Then iterate forward entil end of log.
func (rd *reader) recover() (iter.Seq[*record.Record], error) {
	rcs := list.New()

	for f, err := range rd.sl.Reverse() {
		if err != nil {
			return nil, err
		}

		if done, err := rd.recoverSegment(rcs, f); err != nil {
			return nil, err
		} else if done {
			break
		}
	}

	return iterator.Reverse[*record.Record](rcs), nil
}

func (rd *reader) recoverSegment(rcs *list.List, f *segment.File) (bool, error) {
	size, err := f.Size()
	if err != nil {
		return false, err
	}

	for rc, err := range rd.sr.Reverse(size, f) {
		if err != nil {
			return false, err
		}

		if rc.IsCheckpoint() {
			return true, nil
		}

		rcs.PushBack(rc)
	}

	return false, nil
}

func (rd *reader) reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for f, err := range rd.sl.Reverse() {
			if err != nil {
				yield(nil, err)
				return
			}

			size, err := f.Size()
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range rd.sr.Reverse(size, f) {
				if !yield(rc, err) || err != nil {
					return
				}
			}
		}
	}
}

func (rd *reader) iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for f, err := range rd.sl.IterateFromLSN(lsn) {
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range rd.sr.Iterate(f) {
				if !yield(rc, err) || err != nil {
					return
				}
			}
		}
	}
}

// TODO: Remove this
func (rd *reader) iterateIterator(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return rd.it.Forward(lsn)
}
