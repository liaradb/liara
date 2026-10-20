package fixed

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tip"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
)

// TODO: Create a shared goroutine for each file to manage storage
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

func (fc *FixedCollection) Insert(
	ctx context.Context,
	l span.Log,
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
	w serializer.SizedSerializer,
) error {
	t := tip.NewTip(fc.s, fn)

	s, err := t.Span(ctx, w.Size())
	if err != nil {
		return err
	}

	defer s.Release()

	if err := w.Write(s.Writer(l)); err != nil {
		return err
	}

	s.Commit()

	return fc.c.Insert(ctx, l, fnIdx, k, t.RecordLocator())
}

func (fc *FixedCollection) GetByKey(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
) (*span.Span, error) {
	rl, err := fc.c.Search(ctx, fnIdx, k)
	if err != nil {
		return nil, err
	}

	return fc.GetByRecordLocator(ctx, fn, rl)
}

func (fc *FixedCollection) List(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
) iter.Seq2[*span.Span, error] {
	return func(yield func(*span.Span, error) bool) {
		for rid, err := range fc.c.All(ctx, fnIdx, 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			s, err := fc.GetByRecordLocator(ctx, fn, rid)
			if !yield(s, err) {
				return
			}
		}
	}
}

func (fc *FixedCollection) Replace(
	ctx context.Context,
	l span.Log,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
	k key.Key,
	w serializer.SizedSerializer,
) error {
	s, err := fc.GetByKey(ctx, fn, fnIdx, k)

	defer s.Release()

	if w.Size() > s.Length() {
		// TODO: Verify data can fit
		panic("support resizing")
	}

	if err = w.Write(s.Writer(l)); err != nil {
		return err
	}

	s.Commit()
	return err
}

func (fc *FixedCollection) GetByRecordLocator(
	ctx context.Context,
	fn link.FileName,
	rl link.RecordLocator,
) (*span.Span, error) {
	bid := fn.BlockID(rl.Block())
	b, err := fc.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	s := span.New()
	f := s.AppendSlot(b, rl.SlotID())

	for f.NextPosition() != 0 {
		bid.SetPosition(f.NextPosition())
		b, err := fc.s.Request(ctx, bid)
		if err != nil {
			s.Release()
			return nil, err
		}

		f = s.AppendSlot(b, link.SlotID(0))
	}

	return s, nil
}

func (fc *FixedCollection) Test(
	ctx context.Context,
	fnIdx link.FileName,
	k key.Key,
) (bool, error) {
	_, err := fc.c.Search(ctx, fnIdx, k)
	if errors.Is(err, btree.ErrNotFound) {
		return true, nil
	}

	if err == nil {
		return false, btree.ErrExists
	}

	return false, err
}
