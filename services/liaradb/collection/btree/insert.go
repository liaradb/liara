package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/collection/span"
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
	l span.Log,
	fn link.FileName,
	k key.Key,
	rid link.RecordLocator,
) (err error) {
	var chain *chain
	chain, err = c.getChain(ctx, fn, k)
	if err != nil {
		return err
	}

	defer chain.release()

	chain.latch()
	defer chain.unlatch()

	defer func() {
		if err == nil {
			err = chain.logActions(ctx, l)
		}
	}()

	var fbid link.FileBlockID
	var key = k
	var level byte
	for i, n := range chain.items() {
		var split bool
		if i == 0 {
			ln, ok := n.(*leafnode.LeafNode)
			if !ok {
				err = ErrTypeMismatch
				return
			}

			fbid, key, split, err = c.insertChainLeaf(ctx, fn, fbid, ln, key, rid)
		} else {
			kn, ok := n.(*keynode.KeyNode)
			if !ok {
				err = ErrTypeMismatch
				return
			}

			fbid, key, split, err = c.insertChainKey(ctx, fn, kn, key, fbid.BlockID())
			level++
		}
		if err != nil {
			return
		} else if !split {
			return
		}
	}

	err = c.insertRoot(ctx, fn, level, key, fbid)
	return
}

func (c *insert) getChain(
	ctx context.Context,
	fn link.FileName,
	k key.Key,
) (*chain, error) {
	p, err := c.ns.getPage(ctx, fn.FileBlockID(0))
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

		fbid := fn.FileBlockID(kn.Search(k))
		if p, err = c.ns.getPage(ctx, fbid); err != nil {
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
	fbid link.FileBlockID,
	ln *leafnode.LeafNode,
	k key.Key,
	rid link.RecordLocator,
) (link.FileBlockID, key.Key, bool, error) {
	// Does it exist already?
	if _, ok := ln.Search(k); ok {
		return link.FileBlockID{}, key.Key{}, false, ErrExists
	}

	first, second, ok := ln.Insert(k, rid)
	if ok {
		// no split
		return link.FileBlockID{}, key.Key{}, false, nil
	}

	middle, fbid2, err := c.ns.getNextLeafNode(ctx, fn)
	if err != nil {
		return link.FileBlockID{}, key.Key{}, false, err
	}

	defer middle.Release()

	middle.Latch()
	defer middle.Unlatch()

	rightID := fn.FileBlockID(ln.RightID())
	if rightID.BlockID() != 0 {
		// Only update right node if not root
		right, err := c.ns.getLeafNode(ctx, rightID)
		if err != nil {
			return link.FileBlockID{}, key.Key{}, false, err
		}

		defer right.Release()

		// TODO: Figure out latching
		// ln3.Latch()
		// defer ln3.Unlatch()

		right.SetLeftID(fbid2.BlockID())
	}

	key := middle.Fill(fbid.BlockID(), ln.RightID(), second)
	ln.Replace(fbid2.BlockID(), first)

	return fbid2, key, true, nil
}

// This is a key level page.
//   - Insert, and handle a split.
func (c *insert) insertChainKey(
	ctx context.Context,
	fn link.FileName,
	kn *keynode.KeyNode,
	k key.Key,
	block link.BlockID,
) (link.FileBlockID, key.Key, bool, error) {
	first, second, ok := kn.Insert(k, block)
	if ok {
		return link.FileBlockID{}, key.Key{}, false, nil
	}

	kn2, fbid, err := c.ns.getNextKeyNode(ctx, fn)
	if err != nil {
		return link.FileBlockID{}, key.Key{}, false, err
	}

	defer kn2.Release()

	kn2.Latch()
	defer kn2.Unlatch()

	level := kn.Level()
	key := kn2.Fill(level, second)
	kn.Replace(level, first)

	return fbid, key, true, nil
}

// Created new KeyNode and swap with root
func (c *insert) insertRoot(
	ctx context.Context,
	fn link.FileName,
	level byte,
	key key.Key,
	fbid link.FileBlockID,
) error {
	b0, err := c.ns.getBuffer(ctx, fn.FileBlockID(0))
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
		b2.FileBlockID().BlockID(),
		key,
		fbid.BlockID())

	return nil
}
