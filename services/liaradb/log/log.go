package log

import (
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	size        int64
	segmentList *segment.SegmentList
}

func NewLog(size int64, fsys file.FileSystem, dir string) *Log {
	return &Log{
		size:        size,
		segmentList: segment.NewSegmentList(fsys, dir),
	}
}

func (l *Log) Reader(lsn record.LogSequenceNumber) (*LogReader, error) {
	if err := l.segmentList.Open(); err != nil {
		return nil, err
	}

	_, f, err := l.segmentList.OpenSegmentForLSN(lsn)
	if err != nil {
		return nil, err
	}

	return NewLogReader(l.size, f), nil
}
