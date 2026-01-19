package recovery

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type reader struct {
	sl *segment.List
	sr *segment.Reader
}

func newReader(
	pageSize int64,
	sl *segment.List,
	p page.Page,
) *reader {
	return &reader{
		sl: sl,
		sr: segment.NewReader(pageSize, p),
	}
}

// TODO: Test this
// Iterate in reverse until record type.
//
// Then iterate forward entil end of log.
func (rd *reader) Recover() (iter.Seq[*record.Record], error) {
	rcs := list.New()

	for f, err := range rd.sl.Reverse() {
		if err != nil {
			return nil, err
		}

		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}

		for rc, err := range rd.sr.Reverse(stat.Size(), f) {
			if err != nil {
				return nil, err
			}

			if rc.Action() == record.ActionCheckpoint {
				return listToIterator[*record.Record](rcs), nil
			}

			rcs.PushBack(rc)
		}
	}

	return listToIterator[*record.Record](rcs), nil
}

func (rd *reader) Reverse() iter.Seq2[*record.Record, error] {
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
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(rc, nil) {
					return
				}
			}
		}
	}
}

func (rd *reader) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for f, err := range rd.sl.IterateFromLSN(lsn) {
			if err != nil {
				yield(nil, err)
				return
			}

			for rc, err := range rd.sr.Iterate(f) {
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(rc, nil) {
					return
				}
			}
		}
	}
}

func listToIterator[T any](l *list.List) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value.(T)) {
				return
			}
		}
	}
}
