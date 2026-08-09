package fixedv2

import (
	"context"
	"errors"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/bufferpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
)

type FixedCollection struct {
	s *storage.Storage
	c *btree.Cursor
	l *log.Log
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *FixedCollection {
	return &FixedCollection{
		s: s,
		c: c,
		l: l,
	}
}

func (f *FixedCollection) Insert(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
	v []byte,
) (link.RecordID, error) {
	t := bufferpage.NewTip(f.s, fn)
	s, err := t.Span(ctx, len(v))
	if err != nil {
		return link.RecordID{}, err
	}

	if _, err := s.Write(v); err != nil {
		return link.RecordID{}, err
	}

	s.Commit()

	_, ok := t.Commit()
	if !ok {
		return link.RecordID{}, errors.New("could not commit")
	}

	return t.RecordID(), nil
}
