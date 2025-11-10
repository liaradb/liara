package btree

import (
	"cmp"
	"context"

	"github.com/liaradb/liaradb/storage"
)

type TreeCursor[K cmp.Ordered, V any] struct {
	s       storage.Storage
	current node[K, V]
}

func NewTreeCursor[K cmp.Ordered, V any](s storage.Storage) *TreeCursor[K, V] {
	return &TreeCursor[K, V]{
		s: s,
	}
}

func (tc *TreeCursor[K, V]) GetNode(ctx context.Context, id storage.BlockID) (node[K, V], error) {
	b, err := tc.s.Request(ctx, id)
	if err != nil {
		return nil, err
	}

	bn := &BaseNode[K, V]{
		b: b,
	}

	return bn, nil
}

func (tc *TreeCursor[K, V]) Current() (node[K, V], error) {
	return tc.current, nil
}

func (tc *TreeCursor[K, V]) Left(ctx context.Context) (node[K, V], error) {
	return tc.GetNode(ctx, tc.current.leftID())
}

func (tc *TreeCursor[K, V]) Parent(ctx context.Context) (node[K, V], error) {
	return tc.GetNode(ctx, tc.current.parentID())
}

func (tc *TreeCursor[K, V]) Right(ctx context.Context) (node[K, V], error) {
	return tc.GetNode(ctx, tc.current.rightID())
}
