package pagestorage

import (
	"io"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageStorage struct {
	sl          *segment.List
	wr          filecache.File
	pageSize    int64
	segmentSize action.PageID
	recordSize  int64
	pageID      action.PageID
	timeLineID  action.TimeLineID
}

func New(
	sl *segment.List,
) *PageStorage {
	return &PageStorage{
		sl: sl,
	}
}

func (ps *PageStorage) Init() error {
	_, f, err := ps.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	ps.wr = f
	ps.SeekTail(stat.Size())
	return nil
}

func (ps *PageStorage) SeekTail(size int64) {
	ps.pageID = action.NewActivePageIDFromSize(size, ps.pageSize)
}

func (ps *PageStorage) Append(lsn record.LogSequenceNumber, data []byte) error {
	// TODO: Restore this
	// if err := ps.next(lsn); err != nil {
	// 	return err
	// }

	// return ps.write(data)
	return nil
}

func (ps *PageStorage) Sync(data []byte) error {
	// TODO: Restore this
	// return ps.write(data)
	return nil
}

func (ps *PageStorage) next(lsn record.LogSequenceNumber) error {
	if ps.pageID+1 < ps.segmentSize {
		ps.pageID++
		return nil
	}

	_, f, err := ps.sl.OpenNextSegment(lsn)
	if err != nil {
		return err
	}

	ps.wr = f

	return nil
}

func (ps *PageStorage) position() int64 {
	return ps.pageID.Position(ps.pageSize)
}

func (ps *PageStorage) write(data []byte) error {
	wr := io.NewOffsetWriter(ps.wr, ps.position())
	_, err := wr.Write(data)
	return err
}
