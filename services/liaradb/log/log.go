package log

import "github.com/liaradb/liaradb/file"

type Log struct {
	size        int64
	segmentList *SegmentList
}

func NewLog(size int64, fsys file.FileSystem, dir string) *Log {
	return &Log{
		size:        size,
		segmentList: NewSegmentList(fsys, dir),
	}
}

func (l *Log) Reader(lsn LogSequenceNumber) (*LogReader, error) {
	if err := l.segmentList.Open(); err != nil {
		return nil, err
	}

	_, f, err := l.segmentList.OpenSegmentForLSN(lsn)
	if err != nil {
		return nil, err
	}

	return NewLogReader(l.size, f), nil
}
