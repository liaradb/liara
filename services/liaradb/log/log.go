package log

import (
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	sl     *segment.List
	reader *reader
	writer *writer
}

func NewLog(
	pageSize int64,
	segmentSize page.PageID,
	fsys file.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir)
	return &Log{
		sl:     sl,
		reader: newReader(pageSize, segmentSize, sl),
		writer: newWriter(pageSize, segmentSize, sl),
	}
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.writer.HighWater() }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.writer.LowWater() }
func (l *Log) PageID() page.PageID                 { return l.writer.PageID() }

func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
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
