package pagestorage

import (
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type PageStorage struct {
	sl *segment.List
	f  *segment.File
}

func New(
	sl *segment.List,
) *PageStorage {
	return &PageStorage{
		sl: sl,
	}
}

func (ps *PageStorage) Init() error {
	f, err := ps.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	if err := f.SeekTail(); err != nil {
		return err
	}

	ps.f = f
	return nil
}

func (ps *PageStorage) Append(lsn record.LogSequenceNumber, data []byte) error {
	if err := ps.next(lsn); err != nil {
		return err
	}

	return ps.write(data)
}

func (ps *PageStorage) Sync(data []byte) error {
	return ps.write(data)
}

func (ps *PageStorage) next(lsn record.LogSequenceNumber) error {
	if ps.f.IncrementPageID() {
		return nil
	}

	f, err := ps.sl.OpenNextSegment(lsn)
	if err != nil {
		return err
	}

	ps.f = f

	return nil
}

func (ps *PageStorage) write(data []byte) error {
	return ps.f.Write(data)
}
