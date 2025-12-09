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
		// Leaf
		return c.insertLeaf(page, k, rid)
	} else {
		// Key
		return c.insertKey(ctx, bid.FileName, page, k, rid)
	}
}

func (c *Cursor[K, V]) insertLeaf(
	r page.BTreePage,
	k Key,
	rid RecordID,
) error {
	ln := NewLeafNode(r)
	_, ok := ln.Insert(k, rid)
	if !ok {
		return ErrNoInsert
	}

	return nil
}

func (c *Cursor[K, V]) insertKey(
	ctx context.Context,
	fileName string,
	r page.BTreePage,
	k Key,
	rid RecordID,
) error {
	kn := newKeyNode(r)
	block := kn.Search(k)

	return c.insertPage(ctx,
		storage.NewBlockID(fileName, storage.Offset(block)),
		k,
		rid)
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
