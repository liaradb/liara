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
type insert struct {
	ns *nodeStorage
}

func newInsert(s *storage.Storage) insert {
	return insert{
		ns: newNodeStorage(s),
	}
}

// Insert key value pair into tree
func (c *insert) Insert(
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

			bid, key, split, err = c.insertChainLeaf(ctx, fn, bid, ln, key, rid)
		} else {
			kn, ok := n.(*keynode.KeyNode)
			if !ok {
				return ErrTypeMismatch
			}

			bid, key, split, err = c.insertChainKey(ctx, fn, kn, key, keynode.BlockPosition(bid.Position))
			level++
		}
		if err != nil {
			return err
		} else if !split {
			return nil
		}
	}

	return c.insertRoot(ctx, fn, level, key, bid)
}

func (c *insert) getChain(
	ctx context.Context,
	fn string,
	k key.Key,
) (*chain, error) {
	p, err := c.ns.getPage(ctx, storage.NewBlockID(fn, 0))
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
		if p, err = c.ns.getPage(ctx, bid); err != nil {
			return nil, err
		}
	}

	chain.append(leafnode.New(p))

	return chain, nil
}

// This is a leaf level page.
//   - Insert, and handle a split.
func (c *insert) insertChainLeaf(
	ctx context.Context,
	fn string,
	bid storage.BlockID,
	ln *leafnode.LeafNode,
	k key.Key,
	rid leafnode.RecordID,
) (storage.BlockID, key.Key, bool, error) {
	first, second, ok := ln.Insert(k, rid)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	ln2, bid2, err := c.ns.getNextLeafNode(ctx, fn)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	ln3, err := c.ns.getLeafNode(ctx, storage.NewBlockID(fn, storage.Offset(ln.RightID())))
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer ln2.Release()
	defer ln3.Release()

	ln3.SetLeftID(keynode.BlockPosition(bid2.Position))
	key := ln2.Fill(keynode.BlockPosition(bid.Position), ln.RightID(), second)
	ln.Replace(keynode.BlockPosition(bid2.Position), first)

	return bid2, key, true, nil
}

// This is a key level page.
//   - Insert, and handle a split.
func (c *insert) insertChainKey(
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

	kn2, bid, err := c.ns.getNextKeyNode(ctx, fn)
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
func (c *insert) insertRoot(
	ctx context.Context,
	fn string,
	level byte,
	key key.Key,
	bid storage.BlockID,
) error {
	b0, err := c.ns.getBuffer(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return err
	}

	defer b0.Release()

	// TODO: Should we wrap with KeyNode to simplify latching?
	b2, err := c.ns.getNextBuffer(ctx, fn)
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
