package log

import (
	"io"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type LogWriter struct {
	sw *SegmentWriter
}

func NewLogWriter(
	pageSize int64,
	segmentSize page.PageID,
	rw io.ReadWriteSeeker,
) *LogWriter {
	return &LogWriter{
		sw: NewSegmentWriter(pageSize, segmentSize, rw),
	}
}

func (lw *LogWriter) HighWater() record.LogSequenceNumber { return lw.sw.HighWater() }
func (lw *LogWriter) LowWater() record.LogSequenceNumber  { return lw.sw.LowWater() }
func (lw *LogWriter) PageID() page.PageID                 { return lw.sw.PageID() }

func (lw *LogWriter) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	return lw.sw.Append(rc)
}

func (lw *LogWriter) Flush(lsn record.LogSequenceNumber) error {
	return lw.sw.Flush(lsn)
}

func (lw *LogWriter) Initialize() error {
	return lw.sw.Initialize()
}

func (lw *LogWriter) SeekTail(size int64) error {
	return lw.sw.SeekTail(size)
}
