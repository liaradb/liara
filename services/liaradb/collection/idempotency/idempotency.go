package idempotency

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/record/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type Idempotency struct {
	ic *indexed.IndexedCollection
	c  *btree.Cursor
}

func New(ic *indexed.IndexedCollection, c *btree.Cursor) *Idempotency {
	return &Idempotency{ic: ic, c: c}
}

func (i *Idempotency) Get(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
) (*entity.RequestLog, error) {
	k := key.NewKey(rqid.Bytes())
	s, err := i.ic.GetByKey(ctx, tx, tn, k)
	if err != nil {
		return nil, err
	}

	return i.readRequestLog(s)
}

func (i *Idempotency) List(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
) iter.Seq2[*entity.RequestLog, error] {
	return func(yield func(*entity.RequestLog, error) bool) {
		for s, err := range i.ic.List(ctx, tx, tn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(i.readRequestLog(s)) {
				return
			}
		}
	}
}

func (*Idempotency) readRequestLog(s *span.Span) (*entity.RequestLog, error) {
	defer s.Release()

	rl := entity.RequestLog{}
	return &rl, rl.Read(s)
}

func (i *Idempotency) Insert(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
	e *entity.RequestLog,
) error {
	k := key.NewKey(rqid.Bytes())
	_, err := i.ic.Insert(ctx, tx, tn, k, e)
	return err
}

func (i *Idempotency) Test(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
) (bool, error) {
	k := key.NewKey(rqid.Bytes())
	if err := tx.SLock(ctx, tn, k.String()); err != nil {
		return false, err
	}

	_, err := i.c.Search(ctx, tn.Index(0), k)
	if errors.Is(err, btree.ErrNotFound) {
		return true, nil
	}

	return false, err
}
