package indexed

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/heap/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage/link"
)

// TODO: Create a shared goroutine for each file to manage storage
type IndexedCollection struct {
	h *heap.Heap
	c *btree.Cursor
}

func New(h *heap.Heap, c *btree.Cursor) *IndexedCollection {
	return &IndexedCollection{
		h: h,
		c: c,
	}
}

func (ic *IndexedCollection) Insert(
	ctx context.Context,
	tx Transaction,
	tn tablename.TableName,
	k key.Key,
	w serializer.SizedSerializer,
) (link.RecordID, error) {
	if err := tx.XLock(ctx, tn, k.String()); err != nil {
		return link.RecordID{}, err
	}

	rid, err := ic.h.Insert(ctx, tx, tn.Collection(), w)
	if err != nil {
		return link.RecordID{}, err
	}

	if err := ic.c.Insert(ctx, tx, tn.Index(0), k, rid); err != nil {
		return link.RecordID{}, err
	}

	return rid, nil
}

func (ic *IndexedCollection) GetByKey(
	ctx context.Context,
	tx Transaction,
	tn tablename.TableName,
	k key.Key,
) (*span.Span, error) {
	if err := tx.SLock(ctx, tn, k.String()); err != nil {
		return nil, err
	}

	rid, err := ic.c.Search(ctx, tn.Index(0), k)
	if err != nil {
		return nil, err
	}

	return ic.h.Get(ctx, tn.Collection(), rid)
}

func (ic *IndexedCollection) List(
	ctx context.Context,
	tx Transaction,
	tn tablename.TableName,
) iter.Seq2[*span.Span, error] {
	return func(yield func(*span.Span, error) bool) {
		fn := tn.Collection()
		for r, err := range ic.c.All(ctx, tn.Index(0), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			if err := tx.SLock(ctx, tn, r.Key.String()); err != nil {
				yield(nil, err)
				return
			}

			s, err := ic.h.Get(ctx, fn, r.RecordID)
			if !yield(s, err) {
				return
			}
		}
	}
}

func (ic *IndexedCollection) Replace(
	ctx context.Context,
	tx Transaction,
	tn tablename.TableName,
	k key.Key,
	w serializer.SizedSerializer,
) error {
	if err := tx.XLock(ctx, tn, k.String()); err != nil {
		return err
	}

	rid, err := ic.c.Search(ctx, tn.Index(0), k)
	if err != nil {
		return err
	}

	return ic.h.Replace(ctx, tx, tn.Collection(), rid, w)
}

func (ic *IndexedCollection) GetByRecordID(
	ctx context.Context,
	tn tablename.TableName,
	rid link.RecordID,
) (*span.Span, error) {
	return ic.h.Get(ctx, tn.Collection(), rid)
}
