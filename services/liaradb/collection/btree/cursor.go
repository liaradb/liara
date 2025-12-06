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
