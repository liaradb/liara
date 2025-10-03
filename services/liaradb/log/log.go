package log

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	pageSize    int64
	segmentSize page.PageID
	sl          *segment.SegmentList
	writer      *SegmentWriter
}

func NewLog(
	pageSize int64,
	segmentSize page.PageID,
	fsys file.FileSystem,
	dir string,
) *Log {
	return &Log{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		sl:          segment.NewSegmentList(fsys, dir),
	}
}

func (l *Log) Open() error {
	return l.sl.Open()
}

func (l *Log) Close() error {
	return l.sl.Close()
}

func (l *Log) StartWriter() error {
	_, f, err := l.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	l.writer = NewSegmentWriter(l.pageSize, l.segmentSize, f)
	return l.writer.SeekTail(stat.Size())
}

func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	lsn, err := l.writer.Append(rc)
	if err == ErrInsufficientSpace {
		return l.appendToNextSegment(lsn, rc)
	}

	return lsn, err
}

func (l *Log) appendToNextSegment(lsn record.LogSequenceNumber, rc *record.Record) (record.LogSequenceNumber, error) {
	_, f, err := l.sl.OpenNextSegment(lsn)
	if err != nil {
		return 0, err
	}

	l.writer = NewSegmentWriter(l.pageSize, l.segmentSize, f)
	if err := l.writer.Initialize(); err != nil {
		return 0, err
	}

	return l.writer.Append(rc)
}

func (l *Log) Flush(lsn record.LogSequenceNumber) error {
	return l.writer.Flush(lsn)
}

// TODO: Test this
// Iterate in reverse until record type.
//
// Then iterate forward entil end of log.
func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	r := list.New()

	for f, err := range l.sl.Reverse() {
		if err != nil {
			return nil, err
		}

		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}

		sr := NewSegmentReader(l.pageSize)
		for rc, err := range sr.Reverse(stat.Size(), f) {
			if err != nil {
				return nil, err
			}

			if rc.Action() == record.ActionCheckpoint {
				return listToIterator[*record.Record](r), nil
			}

			r.PushBack(rc)
		}
	}

	return listToIterator[*record.Record](r), nil
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for f, err := range l.sl.Reverse() {
			if err != nil {
				yield(nil, err)
				return
			}

			stat, err := f.Stat()
			if err != nil {
				yield(nil, err)
				return
			}

			sr := NewSegmentReader(l.pageSize)
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

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		for f, err := range l.sl.IterateFromLSN(lsn) {
			if err != nil {
				yield(nil, err)
				return
			}

			sr := NewSegmentReader(l.pageSize)
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
