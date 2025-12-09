package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

type Cursor[K Key, V any] struct {
	s       *storage.Storage
	current node[K, V]
}

func NewCursor[K Key, V any](s *storage.Storage) *Cursor[K, V] {
	return &Cursor[K, V]{
		s: s,
	}
}

func (c *Cursor[K, V]) Insert(
	ctx context.Context,
	fileName string,
	k Key,
	rid RecordID,
) error {
	return c.insertPage(ctx,
		storage.NewBlockID(fileName, 0),
		k,
		rid)
}

func (c *Cursor[K, V]) insertPage(
	ctx context.Context,
	bid storage.BlockID,
	k Key,
	rid RecordID,
) error {
	page, err := c.GetPage(ctx, bid)
	if err != nil {
		return err
	}

	if page.Level() == 0 {
		return c.insertLeaf(page, k, rid)
	} else {
		return c.insertKey(ctx, bid.FileName, page, k, rid)
	}
}

func (c *Cursor[K, V]) insertLeaf(
	p page.BTreePage,
	k Key,
	rid RecordID,
) error {
	ln := NewLeafNode(p)
	_, ok := ln.Insert(k, rid)
	if !ok {
		return ErrNoInsert
	}

	return nil
}

func (c *Cursor[K, V]) insertKey(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
	rid RecordID,
) error {
	kn := newKeyNode(p)
	block := kn.Search(k)

	return c.insertPage(ctx,
		storage.NewBlockID(fileName, storage.Offset(block)),
		k,
		rid)
}

func (c *Cursor[K, V]) Search(
	ctx context.Context,
	fileName string,
	k Key,
) (RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fileName, 0), k)
}

func (c *Cursor[K, V]) searchPage(
	ctx context.Context,
	bid storage.BlockID,
	k Key,
) (RecordID, error) {
	p, err := c.GetPage(ctx, bid)
	if err != nil {
		return RecordID{}, err
	}

	if p.Level() == 0 {
		return c.searchLeaf(p, k)
	} else {
		return c.searchKey(ctx, bid.FileName, p, k)
	}
}

func (*Cursor[K, V]) searchLeaf(
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

func (c *Cursor[K, V]) searchKey(
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

func (c *Cursor[K, V]) GetRoot(ctx context.Context, fileName string) (page.BTreePage, error) {
	return c.GetPage(ctx, storage.NewBlockID(fileName, 0))
}

func (c *Cursor[K, V]) GetPage(ctx context.Context, bid storage.BlockID) (page.BTreePage, error) {
	b, err := c.s.Request(ctx, bid)
	if err != nil {
		return page.BTreePage{}, err
	}

	return page.New(b.Raw()), nil
}

func (c *Cursor[K, V]) GetNode(ctx context.Context, bid storage.BlockID) error {
	bp, err := c.GetPage(ctx, bid)
	if err != nil {
		return err
	}

	if bp.Level() == 0 {
		// Leaf
		_ = NewLeafNode(bp)
	} else {
		// Key
		_ = newKeyNode(bp)
	}

	return nil
}

func (c *Cursor[K, V]) Current() (node[K, V], error) {
	return c.current, nil
}
