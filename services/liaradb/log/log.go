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
	sl          *segment.List
	reader      *LogReader
	writer      *LogWriter
}

func NewLog(
	pageSize int64,
	segmentSize page.PageID,
	fsys file.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir)
	return &Log{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		sl:          sl,
		reader:      NewLogReader(pageSize, segmentSize, sl),
		writer:      NewLogWriter(pageSize, segmentSize, sl),
	}
}

func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	lsn, err := l.writer.Append(rc)
	if err == page.ErrInsufficientSpace {
		return l.appendToNextSegment(lsn, rc)
	}

	return lsn, err
}

func (l *Log) appendToNextSegment(lsn record.LogSequenceNumber, rc *record.Record) (record.LogSequenceNumber, error) {
	_, f, err := l.sl.OpenNextSegment(lsn)
	if err != nil {
		return 0, err
	}

	if err := l.writer.Next(f); err != nil {
		return 0, err
	}

	return l.writer.Append(rc)
}

func (l *Log) Close() error {
	return l.sl.Close()
}

func (l *Log) Flush(lsn record.LogSequenceNumber) error {
	return l.writer.Flush(lsn)
}

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return l.reader.Iterate(lsn)
}

func (l *Log) Open() error {
	return l.sl.Open()
}

func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	return l.reader.Recover()
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return l.reader.Reverse()
}

func (l *Log) StartWriter() error {
	return l.writer.Start()
}
