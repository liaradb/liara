package btreememory

import (
	"cmp"
	"context"
	"errors"
	"slices"

	"github.com/liaradb/liaradb/storage"
)

var id storage.Offset = 0

func nextID() storage.Offset {
	id++
	return id
}

type Storage[K cmp.Ordered] interface {
	GetPage(context.Context) (node[K], error)
	GetRoot(context.Context) (node[K], error)
	SetRoot(context.Context, node[K]) error
}

type Cursor[K cmp.Ordered] struct {
	storage Storage[K]
}

func NewCursor[K cmp.Ordered](s Storage[K]) *Cursor[K] {
	return &Cursor[K]{
		storage: s,
	}
}

func (bt *Cursor[K]) Height(ctx context.Context) (int, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return 0, nil
		}

		return 0, err
	}

	return r.height(), nil
}

func (bt *Cursor[K]) Count(ctx context.Context) (int, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return 0, nil
		}

		return 0, err
	}

	return r.count(), nil
}

func (bt *Cursor[K]) FanOut() int {
	return 3
}

func (bt *Cursor[K]) GetValue(ctx context.Context, k K) (RecordID, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		return RecordID{}, err
	}

	rid, ok := r.getValue(k)
	if !ok {
		return rid, ErrNotFound
	}

	return rid, nil
}

func (bt *Cursor[K]) Insert(ctx context.Context, k K, rid RecordID) error {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			return bt.storage.SetRoot(ctx, newLeafNode(bt.storage, k, rid))
		}

		return err
	}

	var n node[K]
	var ok bool
	switch r := r.(type) {
	case *keyNode[K]:
		n, ok = bt.insertKey(r, k, rid)
	case *leafNode[K]:
		n, ok = bt.insertLeaf(r, k, rid)
	}
	if !ok {
		return ErrNoInsert
	}

	return bt.storage.SetRoot(ctx, newKeyNode(bt.storage, r, n))
}

func (bt *Cursor[K]) insertKey(kn *keyNode[K], k K, rid RecordID) (node[K], bool) {
	n, ok := kn.getChild(k).insert(bt.FanOut(), k, rid)
	if !ok {
		return nil, false
	}

	return kn.insertNode(bt.FanOut(), k, n)
}

func (bt *Cursor[K]) insertLeaf(ln *leafNode[K], k K, rid RecordID) (node[K], bool) {
	c := ln.getChild(k)
	if c != nil {
		// TODO: Create Overflow
		c.value = rid
		return nil, false
	}

	i := ln.getInsertionIndex(k)
	if i == 0 {
		ln.k = k
	}

	// TODO: Split before inserting
	ln.children = slices.Insert(ln.children, i, newLeafEntry(k, rid))
	if len(ln.children) <= bt.FanOut() {
		return nil, false
	}

	return ln.split(), true
}

func (bt *Cursor[K]) DeleteAll(ctx context.Context, k K) error {
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

func (bt *Cursor[K]) DeleteValue(k K, rid RecordID) {

}
