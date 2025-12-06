package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

type TreeCursor[K Key, V any] struct {
	s       storage.Storage
	current node[K, V]
}

func NewTreeCursor[K Key, V any](s storage.Storage) *TreeCursor[K, V] {
	return &TreeCursor[K, V]{
		s: s,
	}
}

func (tc *TreeCursor[K, V]) GetPage(ctx context.Context, bid storage.BlockID) (page.BTreePage, error) {
	b, err := tc.s.Request(ctx, bid)
	if err != nil {
		return page.BTreePage{}, err
	}

	return page.New(b.Raw()), nil
}

func (tc *TreeCursor[K, V]) GetNode(ctx context.Context, bid storage.BlockID) error {
	bp, err := tc.GetPage(ctx, bid)
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

func (tc *TreeCursor[K, V]) Current() (node[K, V], error) {
	return tc.current, nil
}
