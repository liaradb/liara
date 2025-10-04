package log

import (
	"io"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Writer struct {
	highWater     record.LogSequenceNumber
	lowWater      record.LogSequenceNumber
	pageSize      int64
	segmentSize   page.PageID
	sl            *segment.List
	segmentWriter *segment.Writer
}

func NewWriter(
	pageSize int64,
	segmentSize page.PageID,
	sl *segment.List,
) *Writer {
	return &Writer{
		pageSize:    pageSize,
		segmentSize: segmentSize,
		sl:          sl,
	}
}

func (wr *Writer) HighWater() record.LogSequenceNumber { return wr.highWater }
func (wr *Writer) LowWater() record.LogSequenceNumber  { return wr.lowWater }
func (wr *Writer) PageID() page.PageID                 { return wr.segmentWriter.PageID() }

func (wr *Writer) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	lsn, err := wr.appendToSegment(rc)
	if err == page.ErrInsufficientSpace {
		return wr.appendToNextSegment(lsn, rc)
	}

	return lsn, err
}

func (wr *Writer) appendToSegment(rc *record.Record) (record.LogSequenceNumber, error) {
	err := wr.segmentWriter.Append(rc)
	if err != nil {
		if err == page.ErrInsufficientSpace {
			// TODO: Fix this
			return wr.highWater + 1, err
		}
		return 0, err
	}

	wr.highWater++
	return wr.highWater, nil
}

func (wr *Writer) appendToNextSegment(lsn record.LogSequenceNumber, rc *record.Record) (record.LogSequenceNumber, error) {
	_, f, err := wr.sl.OpenNextSegment(lsn)
	if err != nil {
		return 0, err
	}

	if err := wr.Next(f); err != nil {
		return 0, err
	}

	return wr.appendToSegment(rc)
}

func (wr *Writer) Flush(lsn record.LogSequenceNumber) error {
	if err := wr.segmentWriter.Flush(); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, wr.highWater)
	wr.lowWater = lsn
	return nil
}

func (wr *Writer) Initialize() error {
	return wr.segmentWriter.Initialize()
}

func (wr *Writer) SeekTail(size int64) error {
	return wr.segmentWriter.SeekTail(size)
}

func (wr *Writer) Start() error {
	_, f, err := wr.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	wr.segmentWriter = segment.NewWriter(wr.pageSize, wr.segmentSize, f)
	return wr.SeekTail(stat.Size())
}

func (wr *Writer) Next(rw io.ReadWriteSeeker) error {
	wr.segmentWriter = segment.NewWriter(wr.pageSize, wr.segmentSize, rw)
	return wr.Initialize()
}
