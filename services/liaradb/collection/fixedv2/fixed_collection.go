package fixedv2

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/bufferpage"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/span"
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
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
	v []byte,
) error {
	t := bufferpage.NewTip(fc.s, fn)
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
	var bs bufferSlice
	defer bs.Release()

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
	var bs bufferSlice
	defer bs.Release()

	bid := fn.BlockID(rl.Block())
	b, err := fc.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	bs.Append(b)

	p := bufferpage.New(b)
	var s span.Span

	h, d, ok := p.Slot(rl.Position())
	if !ok {
		return nil, errors.New(" could not read slot")
	}
	f := s.Append(h, d)
	for f.Index() != 0 {
		bid = bid.Next()
		b, err := fc.s.Request(ctx, bid)
		if err != nil {
			return nil, err
		}

		bs.Append(b)

		p = bufferpage.New(b)
		h, d, ok := p.Slot(0)
		if !ok {
			return nil, errors.New(" could not read slot")
		}

		f = s.Append(h, d)
	}

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
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
	k key.Key,
	v []byte,
) error {
	var bs bufferSlice
	defer bs.Release()

	rl, err := fc.c.Search(ctx, fnIdx, k)
	if err != nil {
		return err
	}

	bid := fn.BlockID(rl.Block())
	b, err := fc.s.Request(ctx, bid)
	if err != nil {
		return err
	}

	bs.Append(b)

	p := bufferpage.New(b)
	var s span.Span

	h, d, ok := p.Slot(rl.Position())
	if !ok {
		return errors.New(" could not read slot")
	}
	f := s.Append(h, d)
	for f.Index() != 0 {
		bid = bid.Next()
		b, err := fc.s.Request(ctx, bid)
		if err != nil {
			return err
		}

		bs.Append(b)

		p = bufferpage.New(b)
		h, d, ok := p.Slot(0)
		if !ok {
			return errors.New(" could not read slot")
		}

		f = s.Append(h, d)
	}

	_, err = s.Write(v)
	return err
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
