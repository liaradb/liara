package fixed

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/node"
)

// TODO: Create latching
type FixedCollection struct {
	s *storage.Storage
	c *btree.Cursor
}

func New(s *storage.Storage, c *btree.Cursor) *FixedCollection {
	return &FixedCollection{
		s: s,
		c: c,
	}
}

// TODO: Use io.Reader?
func (fc *FixedCollection) Get(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
	k key.Key,
) ([]byte, error) {
	rid, err := fc.c.Search(ctx, fnIdx, k)
	if err != nil {
		return nil, err
	}

	return fc.getItem(ctx, fn, rid)
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

			i, err := fc.getItem(ctx, fn, rid)
			if !yield(i, err) {
				return
			}
		}
	}
}

func (fc *FixedCollection) getItem(ctx context.Context, fn link.FileName, rid link.RecordLocator) ([]byte, error) {
	bid := fn.BlockID(rid.Block())
	b, err := fc.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	defer b.Release()

	n := node.New(b)

	if !n.IsPage() {
		return nil, page.ErrNotPage
	}

	d, ok := n.Child(rid.Position())
	if !ok {
		return nil, btree.ErrNotFound
	}

	return d, nil
}

// TODO: Use io.Writer?
func (fc *FixedCollection) Set(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
	v []byte,
) error {
	crc := page.NewCRC(v)

	rid, ok, err := fc.setCurrent(ctx, fn, v, crc)
	if err != nil {
		return err
	} else if !ok {
		rid, ok, err = fc.setNext(ctx, fn, v, crc)
		if err != nil {
			return err
		} else if !ok {
			return btree.ErrNoInsert
		}
	}

	return fc.c.Insert(ctx, fnIdx, k, rid)
}

func (fc *FixedCollection) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := fc.s.RequestCurrent(ctx, fn)
	if err != nil {
		return link.RecordLocator{}, false, err
	}

	defer b.Release()

	n := node.New(b)
	if !n.IsPage() {
		return link.RecordLocator{}, false, page.ErrNotPage
	}

	rp, d, ok := n.Append(int16(len(v)), crc)
	if !ok {
		return link.RecordLocator{}, false, nil
	}

	copy(d, v)

	return link.NewRecordLocator(b.BlockID().Position(), rp), true, nil
}

func (fc *FixedCollection) setNext(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := fc.s.RequestNext(ctx, fn)
	if err != nil {
		return link.RecordLocator{}, false, err
	}

	defer b.Release()

	n := node.New(b)
	rp, d, ok := n.Append(int16(len(v)), crc)
	if !ok {
		return link.RecordLocator{}, false, nil
	}

	copy(d, v)

	return link.NewRecordLocator(b.BlockID().Position(), rp), true, nil
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
	rid, err := fc.c.Search(ctx, fnIdx, k)
	if err != nil {
		return err
	}

	bid := fn.BlockID(rid.Block())
	b, err := fc.s.Request(ctx, bid)
	if err != nil {
		return err
	}

	defer b.Release()

	n := node.New(b)

	if n.IsPage() {
		return page.ErrNotPage
	}

	if !n.ReplaceChild(int16(rid.Position()), v) {
		return btree.ErrNoUpdate
	}

	return nil
}

func (fc *FixedCollection) Test(
	ctx context.Context,
	fn link.FileName,
	fnIdx link.FileName,
	k key.Key,
) (bool, error) {
	_, err := fc.Get(ctx, fn, fnIdx, value.PartitionID{}, k)
	if errors.Is(err, btree.ErrNotFound) {
		return true, nil
	}

	if err == nil {
		return false, btree.ErrExists
	}

	return false, err
}
