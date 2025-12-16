package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Create latching support
type Cursor struct {
	s *storage.Storage
}

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		s: s,
	}
}

// Insert key value pair into tree
func (c *Cursor) Insert(
	ctx context.Context,
	fn string,
	k Key,
	rid RecordID,
) error {
	chain, err := c.getChain(ctx, fn, k)
	if err != nil {
		return err
	}

	defer chain.release()

	var bid storage.BlockID
	var key = k
	var level byte
	for i, n := range chain.items() {
		var split bool
		if i == 0 {
			ln, ok := n.(*LeafNode)
			if !ok {
				return ErrTypeMismatch
			}

			bid, key, split, err = c.insertChainLeaf(ctx, fn, ln, key, rid)
			if err != nil {
				return err
			} else if !split {
				return nil
			}
		} else {
			kn, ok := n.(*KeyNode)
			if !ok {
				return ErrTypeMismatch
			}

			bid, key, split, err = c.insertChainKey(ctx, fn, kn, key, BlockPosition(bid.Position))
			if err != nil {
				return err
			} else if !split {
				return nil
			}

			level++
		}
	}

	return c.insertRoot(ctx, fn, level, key, bid)
}

func (c *Cursor) getChain(
	ctx context.Context,
	fn string,
	k Key,
) (*chain, error) {
	p, err := c.GetPage(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return nil, err
	}

	chain := newChain()

	for i := int(p.Level()); i > 0; i-- {
		if lvl := p.Level(); lvl != byte(i) {
			return nil, ErrLevelMismatch
		}

		kn := newKeyNode(p)
		chain.append(kn)

		bid := storage.NewBlockID(fn, storage.Offset(kn.Search(k)))
		if p, err = c.GetPage(ctx, bid); err != nil {
			return nil, err
		}
	}

	chain.append(newLeafNode(p))

	return chain, nil
}

// This is a leaf level page.
//   - Insert, and handle a split.
func (c *Cursor) insertChainLeaf(
	ctx context.Context,
	fn string,
	ln *LeafNode,
	k Key,
	rid RecordID,
) (storage.BlockID, Key, bool, error) {
	first, second, ok := ln.Insert(k, rid)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	b, err := c.s.RequestNext(ctx, fn)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer b.Release()

	key := newLeafNode(page.New(b)).Fill(second)
	ln.Replace(first)

	return b.BlockID(), key, true, nil
}

// This is a key level page.
//   - Insert, and handle a split.
func (c *Cursor) insertChainKey(
	ctx context.Context,
	fn string,
	kn *KeyNode,
	k Key,
	block BlockPosition,
) (storage.BlockID, Key, bool, error) {
	first, second, ok := kn.Insert(k, block)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	b, err := c.s.RequestNext(ctx, fn)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer b.Release()

	level := kn.page.Level()
	key := newKeyNode(page.New(b)).Fill(level, second)
	kn.Replace(level, first)

	return b.BlockID(), key, true, nil
}

// Created new KeyNode and swap with root
func (c *Cursor) insertRoot(
	ctx context.Context,
	fn string,
	level byte,
	key Key,
	bid storage.BlockID,
) error {
	b2, err := c.s.RequestNext(ctx, fn)
	if err != nil {
		return err
	}

	defer b2.Release()

	b0, err := c.getBuffer(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return err
	}

	defer b0.Release()

	copy(b2.Raw(), b0.Raw())
	b2.SetDirty()

	root := newKeyNode(page.New(b0))

	// This should always have a child
	child0, _ := root.Child(0)

	// This should always return true
	_ = root.ReplaceRoot(level+1,
		child0.key, BlockPosition(b2.BlockID().Position),
		key, BlockPosition(bid.Position))

	return nil
}

func (c *Cursor) Search(
	ctx context.Context,
	fn string,
	k Key,
) (RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fn, 0), k)
}

func (c *Cursor) searchPage(
	ctx context.Context,
	bid storage.BlockID,
	k Key,
) (RecordID, error) {
	p, err := c.GetPage(ctx, bid)
	if err != nil {
		return RecordID{}, err
	}

	defer p.Release()

	if p.Level() == 0 {
		return c.searchLeaf(p, k)
	} else {
		return c.searchKey(ctx, bid.FileName, p, k)
	}
}

func (*Cursor) searchLeaf(p page.BTreePage, k Key) (RecordID, error) {
	ln := newLeafNode(p)
	rid, ok := ln.Search(k)
	if !ok {
		return RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (c *Cursor) searchKey(
	ctx context.Context,
	fn string,
	p page.BTreePage,
	k Key,
) (RecordID, error) {
	bid := storage.NewBlockID(fn, storage.Offset(newKeyNode(p).Search(k)))
	return c.searchPage(ctx, bid, k)
}

func (c *Cursor) GetPage(ctx context.Context, bid storage.BlockID) (page.BTreePage, error) {
	b, err := c.getBuffer(ctx, bid)
	if err != nil {
		return page.BTreePage{}, err
	}

	return page.New(b), nil
}

func (c *Cursor) getBuffer(ctx context.Context, bid storage.BlockID) (*storage.Buffer, error) {
	return c.s.Request(ctx, bid)
}
