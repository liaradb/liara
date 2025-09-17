package log

import (
	"iter"

	"github.com/liaradb/liaradb/file"
)

type LogReader struct {
	pageSize int64
	file     file.File
}

func NewLogReader(pageSize int64, f file.File) *LogReader {
	return &LogReader{
		pageSize: pageSize,
		file:     f,
	}
}

func (l *LogReader) Iterate() iter.Seq2[*LogRecord, error] {
	return l.IterateFrom(0)
}

func (l *LogReader) IterateFrom(pid LogPageID) iter.Seq2[*LogRecord, error] {
	lpr := newLogPageReader(l.pageSize)
	return lpr.IterateFrom(l.file, pid)
}
