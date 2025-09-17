package log

import (
	"io"
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
	return func(yield func(*LogRecord, error) bool) {
		lpr := newLogPageReader(l.pageSize)
		if err := lpr.Seek(l.file, pid); err != nil {
			yield(nil, err)
			return
		}

		for {
			if _, err := lpr.Read(l.file); err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			for lr, err := range lpr.Records() {
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(lr, nil) {
					return
				}
			}
		}
	}
}
