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
	v []byte,
) error {
	t := tip.NewTip(fc.s, l, fn)
	defer t.Release()

	s, err := t.Span(ctx, len(v))
	if err != nil {
		return err
	}

	if _, err := s.Write(v); err != nil {
		return err
	}

	s.Commit()

	_, ok := t.Commit()
	if !ok {
		return errors.New("could not commit")
	}

	return fc.c.Insert(ctx, fnIdx, k, t.RecordLocator())
}

func (fc *FixedCollection) Get(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
) ([]byte, error) {
	rl, err := fc.c.Search(ctx, fnIdx, k)
	if err != nil {
		return nil, err
	}

	return fc.GetItemByRecordLocator(ctx, fn, rl)
}

func (fc *FixedCollection) List(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for rid, err := range fc.c.All(ctx, fnIdx, 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			i, err := fc.GetItemByRecordLocator(ctx, fn, rid)
			if !yield(i, err) {
				return
			}
		}
	}
}

func (fc *FixedCollection) GetItemByRecordLocator(
	ctx context.Context,
	fn link.FileName,
	rl link.RecordLocator,
) ([]byte, error) {
	// TODO: Don't use nil for Log
	s, err := fc.GetSpanByRecordLocator(ctx, nil, fn, rl)
	if err != nil {
		return nil, err
	}

	defer s.Release()

	// Read Span
	buffer := make([]byte, s.Length())
	if _, err := s.Read(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

// TODO: Use io.Writer?
func (fc *FixedCollection) Replace(
	ctx context.Context,
	l span.Log,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
	k key.Key,
	v []byte,
) error {
	rl, err := fc.c.Search(ctx, fnIdx, k)
	if err != nil {
		return err
	}

	s, err := fc.GetSpanByRecordLocator(ctx, l, fn, rl)
	if err != nil {
		return err
	}

	defer s.Release()

	// TODO: Verify data can fit
	if _, err = s.Write(v); err != nil {
		return err
	}

	if !s.CommitFull() {
		return errors.New("unable to commit")
	}

	return err
}

func (fc *FixedCollection) GetSpanByRecordLocator(
	ctx context.Context,
	l span.Log,
	fn link.FileName,
	rl link.RecordLocator,
) (*span.Span, error) {
	bid := fn.BlockID(rl.Block())
	b, err := fc.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	s := span.New(l)

	f, err := s.AppendSlot(b, rl.SlotID())
	if err != nil {
		b.Release()
		return nil, err
	}

	for f.NextPosition() != 0 {
		bid.SetPosition(f.NextPosition())
		b, err := fc.s.Request(ctx, bid)
		if err != nil {
			s.Release()
			return nil, err
		}

		f, err = s.AppendSlot(b, link.SlotID(0))
		if err != nil {
			b.Release()
			s.Release()
			return nil, err
		}
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
