package log

import (
	"io"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type LogWriter struct {
	highWater     record.LogSequenceNumber
	lowWater      record.LogSequenceNumber
	pageSize      int64
	segmentSize   page.PageID
	sl            *segment.List
	segmentWriter *segment.Writer
}

func NewLogWriter(
	pageSize int64,
	segmentSize page.PageID,
	sl *segment.List,
) *LogWriter {
	return &LogWriter{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		sl:          sl,
	}
}

func (lw *LogWriter) HighWater() record.LogSequenceNumber { return lw.highWater }
func (lw *LogWriter) LowWater() record.LogSequenceNumber  { return lw.lowWater }
func (lw *LogWriter) PageID() page.PageID                 { return lw.segmentWriter.PageID() }

func (lw *LogWriter) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	lsn, err := lw.AppendToSegment(rc)
	if err == page.ErrInsufficientSpace {
		return lw.appendToNextSegment(lsn, rc)
	}

	return lsn, err
}

func (lw *LogWriter) appendToNextSegment(lsn record.LogSequenceNumber, rc *record.Record) (record.LogSequenceNumber, error) {
	_, f, err := lw.sl.OpenNextSegment(lsn)
	if err != nil {
		return 0, err
	}

	if err := lw.Next(f); err != nil {
		return 0, err
	}

	return lw.AppendToSegment(rc)
}

// TODO: Change to private
func (lw *LogWriter) AppendToSegment(rc *record.Record) (record.LogSequenceNumber, error) {
	err := lw.segmentWriter.Append(rc)
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
	if err := lw.segmentWriter.Flush(); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, lw.highWater)
	lw.lowWater = lsn
	return nil
}

func (lw *LogWriter) Initialize() error {
	return lw.segmentWriter.Initialize()
}

func (lw *LogWriter) SeekTail(size int64) error {
	return lw.segmentWriter.SeekTail(size)
}

func (lw *LogWriter) Start() error {
	_, f, err := lw.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	lw.segmentWriter = segment.NewWriter(lw.pageSize, lw.segmentSize, f)
	return lw.SeekTail(stat.Size())
}

func (lw *LogWriter) Next(rw io.ReadWriteSeeker) error {
	lw.segmentWriter = segment.NewWriter(lw.pageSize, lw.segmentSize, rw)
	return lw.Initialize()
}
