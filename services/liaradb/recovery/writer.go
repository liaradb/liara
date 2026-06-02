package recovery

import (
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/segmentio"
)

type writer struct {
	sl *segment.List
	sw *segmentio.Writer
	pq *pagequeue.PageQueue
	ps *pagestorage.PageStorage
}

func newWriter(
	ps *pagestorage.PageStorage,
	pageSize int64,
	segmentSize action.PageID,
	sl *segment.List,
) *writer {
	return &writer{
		sl: sl,
		sw: segmentio.NewWriter(pageSize, segmentSize),
		pq: pagequeue.New(ps, int16(pageSize), page.HeaderSize, page.ItemHeaderSize),
		ps: ps,
	}
}

func (wr *writer) PageID() action.PageID { return wr.sw.PageID() }

func (wr *writer) Append(rc *record.Record) (bool, error) {
	if err := wr.appendToPageQueue(rc); err != nil {
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
func (wr *writer) appendToPageQueue(rc *record.Record) error {
	if err := wr.pq.Append(rc); err != nil {
		return err
	}

	return wr.pq.Flush()
}

func (wr *writer) appendToNextSegment(rc *record.Record, lsn record.LogSequenceNumber) (bool, error) {
	f, err := wr.sl.OpenNextSegment(lsn)
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
	if err := wr.ps.Init(); err != nil {
		return err
	}

	f, err := wr.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	size, err := f.Size()
	if err != nil {
		return err
	}

	return wr.sw.SeekTail(size, f)
}
