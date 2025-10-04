package log

import (
	"iter"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	sl        *segment.List
	reader    *reader
	writer    *writer
	highWater record.LogSequenceNumber
	lowWater  record.LogSequenceNumber
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
		reader: newReader(pageSize, sl),
		writer: newWriter(pageSize, segmentSize, sl),
	}
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.lowWater }
func (l *Log) PageID() page.PageID                 { return l.writer.PageID() }

func (l *Log) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	if err := l.writer.Append(rc, l.highWater+1); err != nil {
		return 0, err
	}

	l.highWater++
	return l.highWater, nil
}

func (l *Log) Close() error {
	return l.sl.Close()
}

func (l *Log) Flush(lsn record.LogSequenceNumber) error {
	if err := l.writer.Flush(lsn); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
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
