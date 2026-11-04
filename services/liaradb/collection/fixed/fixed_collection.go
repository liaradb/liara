package fixed

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage/link"
)

// TODO: Create a shared goroutine for each file to manage storage
// TODO: Separate out Record collection from key value collection
type FixedCollection struct {
	rc *record.RecordCollection
	c  *btree.Cursor
}

func New(rc *record.RecordCollection, c *btree.Cursor) *FixedCollection {
	return &FixedCollection{
		rc: rc,
		c:  c,
	}
}

func (fc *FixedCollection) Insert(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
	w serializer.SizedSerializer,
) (link.RecordID, error) {
	if err := tx.XLock(ctx, tn, k.String()); err != nil {
		return link.RecordID{}, err
	}

	rid, err := fc.rc.Insert(ctx, tx, tn.Collection(), w)
	if err != nil {
		return link.RecordID{}, err
	}

	if err := fc.c.Insert(ctx, tx, tn.Index(0), k, rid); err != nil {
		return link.RecordID{}, err
	}

	return rid, nil
}

func (fc *FixedCollection) GetByKey(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
) (*span.Span, error) {
	if err := tx.SLock(ctx, tn, k.String()); err != nil {
		return nil, err
	}

	rid, err := fc.c.Search(ctx, tn.Index(0), k)
	if err != nil {
		return nil, err
	}

	return fc.rc.Get(ctx, tn.Collection(), rid)
}

func (fc *FixedCollection) List(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
) iter.Seq2[*span.Span, error] {
	return func(yield func(*span.Span, error) bool) {
		fn := tn.Collection()
		for r, err := range fc.c.All(ctx, tn.Index(0), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			if err := tx.SLock(ctx, tn, r.Key.String()); err != nil {
				yield(nil, err)
				return
			}

			s, err := fc.rc.Get(ctx, fn, r.RecordID)
			if !yield(s, err) {
				return
			}
		}
	}
}

func (fc *FixedCollection) Replace(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
	w serializer.SizedSerializer,
) error {
	s, err := fc.getForWritingByKey(ctx, tx, tn, k)
	if err != nil {
		return err
	}

	defer s.Release()

	if w.Size() > s.Length() {
		// TODO: Verify data can fit
		panic("support resizing")
	}

	if err := w.Write(s.Writer()); err != nil {
		return err
	}

	return s.CommitUpdate(ctx, tx)
}

func (fc *FixedCollection) getForWritingByKey(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
) (*span.Span, error) {
	if err := tx.XLock(ctx, tn, k.String()); err != nil {
		return nil, err
	}

	rid, err := fc.c.Search(ctx, tn.Index(0), k)
	if err != nil {
		return nil, err
	}

	return fc.rc.Get(ctx, tn.Collection(), rid)
}

func (fc *FixedCollection) GetByRecordID(
	ctx context.Context,
	tn tablename.TableName,
	rid link.RecordID,
) (*span.Span, error) {
	return fc.rc.Get(ctx, tn.Collection(), rid)
}

func (fc *FixedCollection) Test(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
) (bool, error) {
	if err := tx.SLock(ctx, tn, k.String()); err != nil {
		return false, err
	}

	_, err := fc.c.Search(ctx, tn.Index(0), k)
	if errors.Is(err, btree.ErrNotFound) {
		return true, nil
	}

	if err == nil {
		return false, btree.ErrExists
	}

	return false, err
}
