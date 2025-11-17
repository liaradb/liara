package btreememory

import (
	"cmp"
	"context"
	"errors"
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
			return bt.storage.SetRoot(ctx, newLeafNode(k, v))
		}

		return err
	}

	n, ok := r.insert(bt.FanOut(), k, v)
	if !ok {
		return ErrNoInsert
	}

	return bt.storage.SetRoot(ctx, newKeyNode(r, n))
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
