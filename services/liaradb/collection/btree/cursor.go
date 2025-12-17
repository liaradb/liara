package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Create latching support
// TODO: What happens if two goroutines append simultaneously?
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
	k key.Key,
	rid leafnode.RecordID,
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
			ln, ok := n.(*leafnode.LeafNode)
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
			kn, ok := n.(*keynode.KeyNode)
			if !ok {
				return ErrTypeMismatch
			}

			bid, key, split, err = c.insertChainKey(ctx, fn, kn, key, keynode.BlockPosition(bid.Position))
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
	k key.Key,
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

		kn := keynode.New(p)
		chain.append(kn)

		bid := storage.NewBlockID(fn, storage.Offset(kn.Search(k)))
		if p, err = c.GetPage(ctx, bid); err != nil {
			return nil, err
		}
	}

	chain.append(leafnode.New(p))

	return chain, nil
}

// This is a leaf level page.
//   - Insert, and handle a split.
func (c *Cursor) insertChainLeaf(
	ctx context.Context,
	fn string,
	ln *leafnode.LeafNode,
	k key.Key,
	rid leafnode.RecordID,
) (storage.BlockID, key.Key, bool, error) {
	first, second, ok := ln.Insert(k, rid)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	ln2, bid, err := c.getNextLeafNode(ctx, fn)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer ln2.Release()

	key := ln2.Fill(second)
	ln.Replace(first)

	return bid, key, true, nil
}

// This is a key level page.
//   - Insert, and handle a split.
func (c *Cursor) insertChainKey(
	ctx context.Context,
	fn string,
	kn *keynode.KeyNode,
	k key.Key,
	block keynode.BlockPosition,
) (storage.BlockID, key.Key, bool, error) {
	first, second, ok := kn.Insert(k, block)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	kn2, bid, err := c.getNextKeyNode(ctx, fn)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer kn2.Release()

	level := kn.Level()
	key := kn2.Fill(level, second)
	kn.Replace(level, first)

	return bid, key, true, nil
}

// Created new KeyNode and swap with root
func (c *Cursor) insertRoot(
	ctx context.Context,
	fn string,
	level byte,
	key key.Key,
	bid storage.BlockID,
) error {
	b0, err := c.getBuffer(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return err
	}

	defer b0.Release()

	// TODO: Should we wrap with KeyNode to simplify latching?
	b2, err := c.s.RequestNext(ctx, fn)
	if err != nil {
		return err
	}

	defer b2.Release()

	b2.Clone(b0)

	// This should always return true
	_ = keynode.New(page.New(b0)).ReplaceRoot(
		level+1,
		keynode.BlockPosition(b2.BlockID().Position),
		key,
		keynode.BlockPosition(bid.Position))

	return nil
}

func (c *Cursor) Search(
	ctx context.Context,
	fn string,
	k key.Key,
) (leafnode.RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fn, 0), k)
}

func (c *Cursor) searchPage(
	ctx context.Context,
	bid storage.BlockID,
	k key.Key,
) (leafnode.RecordID, error) {
	p, err := c.GetPage(ctx, bid)
	if err != nil {
		return leafnode.RecordID{}, err
	}

	defer p.Release()

	if p.Level() == 0 {
		return c.searchLeaf(p, k)
	} else {
		return c.searchKey(ctx, bid.FileName, p, k)
	}
}

func (*Cursor) searchLeaf(p page.BTreePage, k key.Key) (leafnode.RecordID, error) {
	ln := leafnode.New(p)
	rid, ok := ln.Search(k)
	if !ok {
		return leafnode.RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (c *Cursor) searchKey(
	ctx context.Context,
	fn string,
	p page.BTreePage,
	k key.Key,
) (leafnode.RecordID, error) {
	bid := storage.NewBlockID(fn, storage.Offset(keynode.New(p).Search(k)))
	return c.searchPage(ctx, bid, k)
}

func (c *Cursor) getNextKeyNode(ctx context.Context, fn string) (*keynode.KeyNode, storage.BlockID, error) {
	b, err := c.s.RequestNext(ctx, fn)
	if err != nil {
		return nil, storage.BlockID{}, err
	}

	return keynode.New(page.New(b)), b.BlockID(), nil
}

func (c *Cursor) getNextLeafNode(ctx context.Context, fn string) (*leafnode.LeafNode, storage.BlockID, error) {
	b, err := c.s.RequestNext(ctx, fn)
	if err != nil {
		return nil, storage.BlockID{}, err
	}

	return leafnode.New(page.New(b)), b.BlockID(), nil
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
