package log

import (
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type Log struct {
	reader *LogReader
}

func NewLog(
	pageSize int64,
	segmentSize page.PageID,
	fsys file.FileSystem,
	dir string,
) *Log {
	return &Log{
		reader: NewLogReader(pageSize, segmentSize, fsys, dir),
	}
}

func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	return l.reader.Append(rc)
}

func (l *Log) Close() error {
	return l.reader.Close()
}

func (l *Log) Flush(lsn record.LogSequenceNumber) error {
	return l.reader.Flush(lsn)
}

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return l.reader.Iterate(lsn)
}

func (l *Log) Open() error {
	return l.reader.Open()
}

func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	return l.reader.Recover()
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return l.reader.Reverse()
}

func (l *Log) StartWriter() error {
	return l.reader.StartWriter()
}
