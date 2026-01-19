package recovery

import (
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type writer struct {
	sl *segment.List
	sw *segment.Writer
}

func newWriter(
	pageSize int64,
	segmentSize page.PageID,
	sl *segment.List,
	p page.Page,
) *writer {
	return &writer{
		sl: sl,
		sw: segment.NewWriter(pageSize, segmentSize, p),
	}
}

func (wr *writer) PageID() page.PageID { return wr.sw.PageID() }

func (wr *writer) Append(rc *record.Record) error {
	err := wr.sw.Append(rc)
	if err == raw.ErrInsufficientSpace {
		err = wr.appendToNextSegment(rc, rc.LogSequenceNumber())
	}

	return err
}

func (wr *writer) appendToNextSegment(rc *record.Record, lsn record.LogSequenceNumber) error {
	_, f, err := wr.sl.OpenNextSegment(lsn)
	if err != nil {
		return err
	}

	wr.sw.Initialize(f)

	return wr.sw.Append(rc)
}

func (wr *writer) Flush(lsn record.LogSequenceNumber) error {
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
