package btreememory

import (
	"cmp"
	"context"
	"errors"
	"slices"
)

type Storage[K cmp.Ordered, V any] interface {
	GetPage(context.Context) (node[K, V], error)
	GetRoot(context.Context) (node[K, V], error)
	SetRoot(context.Context, node[K, V]) error
}

type Cursor[K cmp.Ordered, V any] struct {
	storage Storage[K, V]
}

func NewCursor[K cmp.Ordered, V any](s Storage[K, V]) *Cursor[K, V] {
	return &Cursor[K, V]{
		storage: s,
	}
}

func (bt *Cursor[K, V]) Height(ctx context.Context) (int, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return 0, nil
		}

		return 0, err
	}

	return r.height(), nil
}

func (bt *Cursor[K, V]) Count(ctx context.Context) (int, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return 0, nil
		}

		return 0, err
	}

	return r.count(), nil
}

func (bt *Cursor[K, V]) FanOut() int {
	return 3
}

func (bt *Cursor[K, V]) GetValue(ctx context.Context, k K) (V, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		var v V
		return v, err
	}

	v, ok := r.getValue(k)
	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (bt *Cursor[K, V]) Insert(ctx context.Context, k K, v V) error {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return bt.storage.SetRoot(ctx, newLeafNode(bt.storage, k, v))
		}

		return err
	}

	var n node[K, V]
	var ok bool
	switch r := r.(type) {
	case *keyNode[K, V]:
		n, ok = bt.insertKey(r, k, v)
	case *leafNode[K, V]:
		n, ok = bt.insertLeaf(r, k, v)
	}
	if !ok {
		return ErrNoInsert
	}

	return bt.storage.SetRoot(ctx, newKeyNode(bt.storage, r, n))
}

func (bt *Cursor[K, V]) insertKey(kn *keyNode[K, V], k K, v V) (node[K, V], bool) {
	n, ok := kn.getChild(k).insert(bt.FanOut(), k, v)
	if !ok {
		return nil, false
	}

	return kn.insertNode(bt.FanOut(), k, n)
}

func (bt *Cursor[K, V]) insertLeaf(ln *leafNode[K, V], k K, v V) (node[K, V], bool) {
	c := ln.getChild(k)
	if c != nil {
		// TODO: Create Overflow
		c.append(v)
		return nil, false
	}

	i := ln.getInsertionIndex(k)
	if i == 0 {
		ln.k = k
	}

	// TODO: Split before inserting
	ln.children = slices.Insert(ln.children, i, newLeafEntry(k, v))
	if len(ln.children) <= bt.FanOut() {
		return nil, false
	}

	return ln.split(), true
}

func (bt *Cursor[K, V]) DeleteAll(ctx context.Context, k K) error {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return nil
		}

		return err
	}

	r.deleteAll(bt.FanOut(), k)
	return nil
}

func (bt *Cursor[K, V]) DeleteValue(k K, v V) {

}
