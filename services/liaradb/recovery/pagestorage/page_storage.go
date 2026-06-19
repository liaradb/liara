package pagestorage

import (
	"errors"

	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

// TODO: Should we combine this into segment.List?
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

func (ps *PageStorage) Init(data []byte) error {
	f, err := ps.sl.OpenLatestSegment()
	if err != nil {
		return err
	}

	size, err := f.SeekTail()
	if err != nil {
		return errors.Join(err, f.Close())
	}

	if size == 0 {
		clear(data)
	} else {
		if _, err := f.Read(data); err != nil {
			return errors.Join(err, f.Close())
		}
	}

	ps.f = f
	return nil
}

func (ps *PageStorage) Append(lsn record.LogSequenceNumber, data []byte) error {
	if err := ps.nextPage(lsn); err != nil {
		return err
	}

	return ps.write(data)
}

func (ps *PageStorage) nextPage(lsn record.LogSequenceNumber) error {
	if ps.f.NextPage() {
		return nil
	}

	return ps.nextSegment(lsn)
}

func (ps *PageStorage) nextSegment(lsn record.LogSequenceNumber) error {
	f, err := ps.sl.OpenNextSegment(lsn)
	if err != nil {
		return err
	}

	ps.f = f
	return nil
}

func (ps *PageStorage) Sync(data []byte) error {
	return ps.write(data)
}

func (ps *PageStorage) write(data []byte) error {
	return ps.f.Write(data)
}
