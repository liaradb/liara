package recovery

import (
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type writer struct {
	sl *segment.List
	sw *segment.Writer
}

func newWriter(
	pageSize int64,
	segmentSize action.PageID,
	recordSize int64,
	sl *segment.List,
) *writer {
	return &writer{
		sl: sl,
		sw: segment.NewWriter(pageSize, segmentSize, recordSize),
	}
}

func (wr *writer) PageID() action.PageID { return wr.sw.PageID() }

func (wr *writer) Append(rc *record.Record) (bool, error) {
	flushed, err := wr.sw.Append(rc)
	if err == raw.ErrInsufficientSpace {
		// Ignore this flushed value, as it's the first record
		_, err = wr.appendToNextSegment(rc, rc.LogSequenceNumber())
	}

	return flushed, err
}

func (wr *writer) appendToNextSegment(rc *record.Record, lsn record.LogSequenceNumber) (bool, error) {
	_, f, err := wr.sl.OpenNextSegment(lsn)
	if err != nil {
		return false, err
	}

	wr.sw.Initialize(f)
	return wr.sw.Append(rc)
}

func (wr *writer) Flush() error {
	return wr.sw.Flush()
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

	return wr.sw.SeekTail(stat.Size(), f)
}
