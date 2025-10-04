package log

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type reader struct {
	pageSize    int64
	segmentSize page.PageID
	sl          *segment.List
}

func newReader(
	pageSize int64,
	segmentSize page.PageID,
	sl *segment.List,
) *reader {
	return &reader{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		sl:          sl,
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

		sr := segment.NewReader(rd.pageSize)
		for rc, err := range sr.Reverse(stat.Size(), f) {
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

			sr := segment.NewReader(rd.pageSize)
			for rc, err := range sr.Reverse(stat.Size(), f) {
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

			sr := segment.NewReader(rd.pageSize)
			for rc, err := range sr.Iterate(f) {
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
