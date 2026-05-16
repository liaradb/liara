package pagestorage

import (
	"io"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageStorage struct {
	sl          *segment.List
	wr          io.WriterAt
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

	ps.wr = f

	return nil
}

func (ps *PageStorage) Append([]byte) error {
	return nil
}

func (ps *PageStorage) Sync([]byte) error {
	return nil
}
