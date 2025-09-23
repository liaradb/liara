package log

import (
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	size int64
	sl   *segment.SegmentList
}

func NewLog(size int64, fsys file.FileSystem, dir string) *Log {
	return &Log{
		size: size,
		sl:   segment.NewSegmentList(fsys, dir),
	}
}

func (l *Log) Open() error {
	return l.sl.Open()
}

func (l *Log) Close() error {
	return l.sl.Close()
}

func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	_, f, err := l.sl.OpenLatestSegment()
	if err != nil {
		return 0, err
	}

	lw := NewLogWriter(l.size, f)
	lsn, err := lw.Append(rc)
	if err != nil {
		return 0, err
	}

	if err := lw.Flush(lsn); err != nil {
		return 0, err
	}

	return lsn, nil
}

func (l *Log) Recover() error {
	// TODO: Implement this
	panic("unimplemented")
}

// TODO: Create SegmentList reverse iterator
func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	// TODO: Implement this
	panic("unimplemented")
}

// TODO: Create SegmentList iterator
func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		lr, err := l.reader(lsn)
		if err != nil {
			yield(nil, err)
			return
		}

		for rc, err := range lr.Iterate() {
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

func (l *Log) reader(lsn record.LogSequenceNumber) (*LogReader, error) {
	_, f, err := l.sl.OpenSegmentForLSN(lsn)
	if err != nil {
		return nil, err
	}

	return NewLogReader(l.size, f), nil
}
