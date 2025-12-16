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

func (c *Cursor) getChain(
	ctx context.Context,
	fileName string,
	k Key,
) (*chain, error) {
	bid := storage.NewBlockID(fileName, 0)
	p, err := c.GetPage(ctx, bid)
	if err != nil {
		return nil, err
	}

	chain := newChain()

	level := p.Level()
	if level == 0 {
		// leaf
		ln := NewLeafNode(p)
		chain.append(ln)
		return chain, nil
	}

	for i := int(level); i >= 0; i-- {
		lvl := p.Level()
		if lvl != byte(i) {
			return nil, ErrLevelMismatch
		}

		if i == 0 {
			// leaf
			ln := NewLeafNode(p)
			chain.append(ln)
			break
		}

		kn := newKeyNode(p)
		chain.append(kn)
		block := kn.Search(k)

		if p, err = c.GetPage(ctx, storage.NewBlockID(fileName, storage.Offset(block))); err != nil {
			return nil, err
		}
	}

	return chain, nil
}

// Insert key value pair into tree
func (c *Cursor) Insert(
	ctx context.Context,
	fileName string,
	k Key,
	rid RecordID,
) error {
	chain, err := c.getChain(ctx, fileName, k)
	if err != nil {
		return err
	}

	defer chain.release()

	var bid storage.BlockID
	var key = k
	var level byte
	var split bool
	for i, n := range chain.items() {
		if i == 0 {
			ln, ok := n.(*LeafNode)
			if !ok {
				return ErrTypeMismatch
			}

			bid, key, split, err = c.insertChainLeaf(ctx, fileName, ln, key, rid)
			if err != nil {
				return err
			}
		} else {
			// TODO: Split Key Node
			if !split {
				return nil
			}

			kn, ok := n.(*KeyNode)
			if !ok {
				return ErrTypeMismatch
			}

			bid, key, split, err = c.insertChainKey(ctx, fileName, kn, key, BlockPosition(bid.Position))
			if err != nil {
				return err
			}
			level++
		}
		i++
	}

	if !split {
		return nil
	}

	bid0 := storage.NewBlockID(fileName, 0)

	// Swap block2 with root
	b2, err := c.s.RequestNext(ctx, fileName)
	if err != nil {
		return err
	}

	defer b2.Release()

	b0, err := c.GetBuffer(ctx, bid0)
	if err != nil {
		return err
	}

	defer b0.Release()

	root := newKeyNode(page.New(b0))

	// This should always have a child
	child0, _ := root.Child(0)

	copy(b2.Raw(), b0.Raw())
	b2.SetDirty()

	root.page.Clear()
	root.page.SetLevel(level + 1)
	_, _ = root.Append(child0.key, BlockPosition(b2.BlockID().Position))
	_, _ = root.Append(key, BlockPosition(bid.Position))

	return nil
}

// This is a leaf level page.
//   - Insert, and handle a split.
func (c *Cursor) insertChainLeaf(
	ctx context.Context,
	fileName string,
	ln *LeafNode,
	k Key,
	rid RecordID,
) (storage.BlockID, Key, bool, error) {
	first, second, ok := ln.Insert(k, rid)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	b, err := c.s.RequestNext(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer b.Release()

	p2 := page.New(b)
	ln2 := NewLeafNode(p2)

	key := ln2.Fill(second)
	ln.Replace(first)

	return b.BlockID(), key, true, nil
}

// This is a key level page.
//   - Insert, and handle a split.
func (c *Cursor) insertChainKey(
	ctx context.Context,
	fileName string,
	kn *KeyNode,
	k Key,
	block BlockPosition,
) (storage.BlockID, Key, bool, error) {
	first, second, ok := kn.Insert(k, block)
	if ok {
		return storage.BlockID{}, "", false, nil
	}

	b, err := c.s.RequestNext(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, "", false, err
	}

	defer b.Release()

	level := kn.page.Level()

	p2 := page.New(b)
	kn2 := newKeyNode(p2)
	kn2.page.SetLevel(level)

	key := kn2.Fill(second)
	kn.Replace(first)
	// TODO: Clean this
	kn.page.SetLevel(level)

	return b.BlockID(), key, true, nil
}

func (c *Cursor) Search(
	ctx context.Context,
	fileName string,
	k Key,
) (RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fileName, 0), k)
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

func (*Cursor) searchLeaf(
	p page.BTreePage,
	k Key,
) (RecordID, error) {
	ln := NewLeafNode(p)
	rid, ok := ln.Search(k)
	if !ok {
		return RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (c *Cursor) searchKey(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
) (RecordID, error) {
	kn := newKeyNode(p)
	block := kn.Search(k)
	return c.searchPage(ctx,
		storage.NewBlockID(fileName, storage.Offset(block)),
		k)
}

func (c *Cursor) GetRoot(ctx context.Context, fileName string) (page.BTreePage, error) {
	return c.GetPage(ctx, storage.NewBlockID(fileName, 0))
}

func (c *Cursor) GetPage(ctx context.Context, bid storage.BlockID) (page.BTreePage, error) {
	b, err := c.GetBuffer(ctx, bid)
	if err != nil {
		return page.BTreePage{}, err
	}

	return page.New(b), nil
}

func (c *Cursor) GetBuffer(ctx context.Context, bid storage.BlockID) (*storage.Buffer, error) {
	return c.s.Request(ctx, bid)
}
