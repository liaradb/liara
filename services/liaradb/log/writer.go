package log

import (
	"io"

	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type writer struct {
	highWater     record.LogSequenceNumber
	lowWater      record.LogSequenceNumber
	sl            *segment.List
	segmentWriter *segment.Writer
}

func newWriter(
	pageSize int64,
	segmentSize page.PageID,
	sl *segment.List,
) *writer {
	return &writer{
		sl:            sl,
		segmentWriter: segment.NewWriter(pageSize, segmentSize),
	}
}

func (wr *writer) HighWater() record.LogSequenceNumber { return wr.highWater }
func (wr *writer) LowWater() record.LogSequenceNumber  { return wr.lowWater }
func (wr *writer) PageID() page.PageID                 { return wr.segmentWriter.PageID() }

func (wr *writer) Append(rc *record.Record) (record.LogSequenceNumber, error) {
	err := wr.appendToSegment(rc)
	if err == page.ErrInsufficientSpace {
		err = wr.appendToNextSegment(rc)
	}

	if err != nil {
		return 0, err
	}

	wr.highWater++
	return wr.highWater, nil
}

func (wr *writer) appendToSegment(rc *record.Record) error {
	err := wr.segmentWriter.Append(rc)
	if err != nil {
		if err == page.ErrInsufficientSpace {
			return err
		}
		return err
	}

	return nil
}

func (wr *writer) appendToNextSegment(rc *record.Record) error {
	_, f, err := wr.sl.OpenNextSegment(wr.highWater)
	if err != nil {
		return err
	}

	if err := wr.next(f); err != nil {
		return err
	}

	return wr.appendToSegment(rc)
}

func (wr *writer) next(rw io.ReadWriteSeeker) error {
	return wr.segmentWriter.Initialize(rw)
}

func (wr *writer) Flush(lsn record.LogSequenceNumber) error {
	if err := wr.segmentWriter.Flush(); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = min(lsn, wr.highWater)
	wr.lowWater = lsn
	return nil
}

func (wr *writer) Start() error {
	_, f, err := wr.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	return wr.seekTail(stat.Size(), f)
}

func (wr *writer) seekTail(size int64, rw io.ReadWriteSeeker) error {
	return wr.segmentWriter.SeekTail(size, rw)
}
