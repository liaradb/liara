package recovery

import (
	encoder "github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/segmentio"
)

type writer struct {
	sl *segment.List
	sw *segmentio.Writer
	pq *pagequeue.PageQueue
}

func newWriter(
	ps pagequeue.PageStorage,
	pageSize int64,
	segmentSize action.PageID,
	recordSize int64,
	sl *segment.List,
) *writer {
	return &writer{
		sl: sl,
		sw: segmentio.NewWriter(pageSize, segmentSize, recordSize),
		pq: pagequeue.New(ps, int16(pageSize), page.HeaderSize, page.ItemHeaderSize),
	}
}

func (wr *writer) PageID() action.PageID { return wr.sw.PageID() }
func (wr *writer) RecordSize() int64     { return wr.sw.RecordSize() }

func (wr *writer) Append(rc *record.Record) (bool, error) {
	if _, err := wr.appendToPageQueue(rc); err != nil {
		return false, err
	}

	flushed, err := wr.sw.AppendRecord(rc)
	if err == raw.ErrInsufficientSpace {
		// Ignore this flushed value, as it's the first record
		_, err = wr.appendToNextSegment(rc, rc.LogSequenceNumber())
	}

	return flushed, err
}

// # Append to PageQueue
//   - If current page was full already, do not sync
//   - Otherwise, sync current page
//   - For every page after, push new page to disk
//   - For last page, store that as current
func (wr *writer) appendToPageQueue(rc *record.Record) (*encoder.Page, error) {
	if err := wr.pq.Append(rc); err != nil {
		return nil, err
	}

	var current *encoder.Page
	for p := range wr.pq.Pages() {
		_ = p.Data()
		current = p
	}

	if err := wr.pq.Flush(); err != nil {
		return nil, err
	}

	return current, nil
}

func (wr *writer) appendToNextSegment(rc *record.Record, lsn record.LogSequenceNumber) (bool, error) {
	_, f, err := wr.sl.OpenNextSegment(lsn)
	if err != nil {
		return false, err
	}

	wr.sw.Initialize(f)
	return wr.sw.AppendRecord(rc)
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
