package idempotency

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage"
)

type Idempotency struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor) *Idempotency {
	return &Idempotency{
		fc: fixed.New(record.NewRecordCollection(s), c),
	}
}

func (i *Idempotency) Get(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
) (*entity.RequestLog, error) {
	k := key.NewKey(rqid.Bytes())
	s, err := i.fc.GetByKey(ctx, tx, tn, k)
	if err != nil {
		return nil, err
	}

	return i.readRequestLog(s)
}

func (i *Idempotency) List(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
) iter.Seq2[*entity.RequestLog, error] {
	return func(yield func(*entity.RequestLog, error) bool) {
		for s, err := range i.fc.List(ctx, tx, tn) {
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
	tx fixed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
	e *entity.RequestLog,
) error {
	k := key.NewKey(rqid.Bytes())
	_, err := i.fc.Insert(ctx, tx, tn, k, e)
	return err
}

func (i *Idempotency) Test(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
) (bool, error) {
	k := key.NewKey(rqid.Bytes())
	return i.fc.Test(ctx, tx, tn, k)
}
