package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

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
	fn link.FileName,
	k value.Key,
	rid link.RecordLocator,
) error {
	chain, err := c.getChain(ctx, fn, k)
	if err != nil {
		return err
	}

	defer chain.release()

	chain.latch()
	defer chain.unlatch()

	var bid link.BlockID
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

			bid, key, split, err = c.insertChainKey(ctx, fn, kn, key, bid.Position())
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
	fn link.FileName,
	k value.Key,
) (*chain, error) {
	p, err := c.ns.getPage(ctx, fn.BlockID(0))
	if err != nil {
		return nil, err
	}

	chain := newChain()

	for i := int(p.Level()); i > 0; i-- {
		if lvl := p.Level(); lvl != byte(i) {
			chain.release()
			return nil, ErrLevelMismatch
		}

		kn := keynode.New(p)
		chain.append(kn)

		bid := fn.BlockID(kn.Search(k))
		if p, err = c.ns.getPage(ctx, bid); err != nil {
			chain.release()
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
	fn link.FileName,
	bid link.BlockID,
	ln *leafnode.LeafNode,
	k value.Key,
	rid link.RecordLocator,
) (link.BlockID, value.Key, bool, error) {
	first, second, ok := ln.Insert(k, rid)
	if ok {
		return link.BlockID{}, value.Key{}, false, nil
	}

	ln2, bid2, err := c.ns.getNextLeafNode(ctx, fn)
	if err != nil {
		return link.BlockID{}, value.Key{}, false, err
	}

	defer ln2.Release()

	ln2.Latch()
	defer ln2.Unlatch()

	ln3, err := c.ns.getLeafNode(ctx, fn.BlockID(ln.RightID()))
	if err != nil {
		return link.BlockID{}, value.Key{}, false, err
	}

	defer ln3.Release()

	// TODO: Figure out latching
	// ln3.Latch()
	// defer ln3.Unlatch()

	ln3.SetLeftID(bid2.Position())
	key := ln2.Fill(bid.Position(), ln.RightID(), second)
	ln.Replace(bid2.Position(), first)

	return bid2, key, true, nil
}

// This is a key level page.
//   - Insert, and handle a split.
func (c *insert) insertChainKey(
	ctx context.Context,
	fn link.FileName,
	kn *keynode.KeyNode,
	k value.Key,
	block link.FilePosition,
) (link.BlockID, value.Key, bool, error) {
	first, second, ok := kn.Insert(k, block)
	if ok {
		return link.BlockID{}, value.Key{}, false, nil
	}

	kn2, bid, err := c.ns.getNextKeyNode(ctx, fn)
	if err != nil {
		return link.BlockID{}, value.Key{}, false, err
	}

	defer kn2.Release()

	kn2.Latch()
	defer kn2.Unlatch()

	level := kn.Level()
	key := kn2.Fill(level, second)
	kn.Replace(level, first)

	return bid, key, true, nil
}

// Created new KeyNode and swap with root
func (c *insert) insertRoot(
	ctx context.Context,
	fn link.FileName,
	level byte,
	key value.Key,
	bid link.BlockID,
) error {
	b0, err := c.ns.getBuffer(ctx, fn.BlockID(0))
	if err != nil {
		return err
	}

	defer b0.Release()

	// TODO: Figure out latching
	// b0.Latch()
	// defer b0.Unlatch()

	// TODO: Should we wrap with KeyNode to simplify latching?
	b2, err := c.ns.getNextBuffer(ctx, fn)
	if err != nil {
		return err
	}

	defer b2.Release()

	b2.Latch()
	defer b2.Unlatch()

	b2.Clone(b0)

	// This should always return true
	_ = keynode.New(node.New(b0)).ReplaceRoot(
		level+1,
		b2.BlockID().Position(),
		key,
		bid.Position())

	return nil
}
