package log

import (
	"io"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

type LogWriter struct {
	highWater record.LogSequenceNumber
	lowWater  record.LogSequenceNumber
	sw        *SegmentWriter
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

func (lw *LogWriter) HighWater() record.LogSequenceNumber { return lw.highWater }
func (lw *LogWriter) LowWater() record.LogSequenceNumber  { return lw.lowWater }
func (lw *LogWriter) PageID() page.PageID                 { return lw.sw.PageID() }

func (lw *LogWriter) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	err := lw.sw.Append(rc)
	if err != nil {
		if err == page.ErrInsufficientSpace {
			// TODO: Fix this
			return lw.highWater + 1, err
		}
		return 0, err
	}

	lw.highWater++
	return lw.highWater, nil
}

func (lw *LogWriter) Flush(lsn record.LogSequenceNumber) error {
	if err := lw.sw.Flush(); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, lw.highWater)
	lw.lowWater = lsn
	return nil
}

func (lw *LogWriter) Initialize() error {
	return lw.sw.Initialize()
}

func (lw *LogWriter) SeekTail(size int64) error {
	return lw.sw.SeekTail(size)
}
