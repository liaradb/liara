package log

import (
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
	writer      *LogWriter
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

	l.writer = NewLogWriter(l.pageSize, l.segmentSize, f)
	return nil
}

// TODO: Clean this
func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	lsn, err := l.writer.Append(rc)
	if err != nil {
		if err == ErrInsufficientSpace {
			_, f, err := l.sl.OpenNextSegment(lsn)
			if err != nil {
				return 0, err
			}
			l.writer = NewLogWriter(l.pageSize, l.segmentSize, f)
			return l.writer.Append(rc)
		}
	}

	return lsn, nil
}

func (l *Log) Flush(lsn record.LogSequenceNumber) error {
	return l.writer.Flush(lsn)
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
		for f, err := range l.sl.IterateFromLSN(lsn) {
			if err != nil {
				yield(nil, err)
				return
			}

			lr := NewLogReader(l.pageSize, f)
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
}
