package btreememory

import (
	"cmp"
	"context"
	"errors"

	"github.com/liaradb/liaradb/storage"
)

var id storage.Offset = -1

func nextID() storage.Offset {
	id++
	return id
}

type Storage[K cmp.Ordered] interface {
	GetNode(context.Context, storage.BlockID) (node[K], error)
	GetKeyNode(context.Context, storage.BlockID) (*keyNode[K], error)
	GetLeafNode(context.Context, storage.BlockID) (*leafNode[K], error)
	InsertNode(context.Context, storage.BlockID, node[K]) error
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

// func (bt *Cursor[K]) Count(ctx context.Context) (int, error) {
// 	r, err := bt.storage.GetRoot(ctx)
// 	if err != nil {
// 		if errors.Is(err, ErrEmptyTree) {
// 			return 0, nil
// 		}

// 		return 0, err
// 	}

// 	return r.count(), nil
// }

func (bt *Cursor[K]) FanOut() int {
	return 3
}

func (bt *Cursor[K]) GetValue(ctx context.Context, k K) (RecordID, error) {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		return RecordID{}, err
	}

	if off, ok := r.getChild(k); ok {
		return bt.getChild(ctx, k, off)
	}

	if rid, ok := r.getValue(k); ok {
		return rid, nil
	} else {
		return RecordID{}, ErrNotFound
	}
}

func (bt *Cursor[K]) getChild(ctx context.Context, k K, off storage.Offset) (RecordID, error) {
	n, err := bt.storage.GetNode(ctx, storage.NewBlockID("", off))
	if err != nil {
		return RecordID{}, err
	}

	if off, ok := n.getChild(k); ok {
		return bt.getChild(ctx, k, off)
	}

	if rid, ok := n.getValue(k); ok {
		return rid, nil
	} else {
		return RecordID{}, ErrNotFound
	}
}

func (bt *Cursor[K]) Insert(ctx context.Context, k K, rid RecordID) error {
	r, err := bt.storage.GetRoot(ctx)
	if err != nil {
		if errors.Is(err, ErrEmptyTree) {
			ln := newLeafNode(bt.storage, k, rid)
			if err := bt.storage.InsertNode(ctx, storage.NewBlockID("", ln.i), ln); err != nil {
				return err
			}

			return bt.storage.SetRoot(ctx, ln)
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

	return bt.storage.SetRoot(ctx, n)
}

func (bt *Cursor[K]) insertKey(kn *keyNode[K], k K, rid RecordID) (node[K], bool) {
	n, ok := kn.getChild(k)
	if !ok {
		return nil, false
	}

	if kn.level == 2 {
		// Child is a leafNode
		ln, _ := bt.storage.GetLeafNode(context.Background(), storage.NewBlockID("", n))
		ln2, ok := ln.insert(bt.FanOut(), k, rid)
		if !ok {
			return kn, true
		}

		kn2 := newKeyNode(bt.storage, ln, ln2)
		_ = bt.storage.InsertNode(context.Background(), storage.NewBlockID("", kn2.i), kn2)
		return kn2, true
	} else {
		// Child is a keyNode
	}

	// if !ok {
	// 	return nil, false
	// }

	// .insert(bt.FanOut(), k, rid)

	return kn.insert(bt.FanOut(), k, n)
}

func (bt *Cursor[K]) insertLeaf(ln *leafNode[K], k K, rid RecordID) (node[K], bool) {
	ln2, ok := ln.insert(bt.FanOut(), k, rid)
	if !ok {
		return ln, true
	}

	_ = bt.storage.InsertNode(context.Background(), storage.NewBlockID("", ln2.i), ln2)

	kn2 := newKeyNode(bt.storage, ln, ln2)
	_ = bt.storage.InsertNode(context.Background(), storage.NewBlockID("", kn2.i), kn2)
	return kn2, true
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
