package recovery

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/util/iterator"
)

type reader struct {
	sl *segment.List
	sr *segment.Reader
}

func newReader(
	pageSize int64,
	sl *segment.List,
) *reader {
	return &reader{
		sl: sl,
		sr: segment.NewReader(pageSize),
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

func (rd *reader) recoverSegment(rcs *list.List, f filecache.File) (bool, error) {
	stat, err := f.Stat()
	if err != nil {
		return false, err
	}

	for rc, err := range rd.sr.Reverse(stat.Size(), f) {
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

			stat, err := f.Stat()
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range rd.sr.Reverse(stat.Size(), f) {
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
