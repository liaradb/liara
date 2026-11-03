package outbox

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/log"
)

type Outbox struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *Outbox {
	return &Outbox{
		fc: fixed.New(s, c, l),
	}
}

func (o *Outbox) Get(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	k := key.NewKey(oid.Bytes())
	s, err := o.fc.GetByKey(ctx, tx, tn, k)
	if err != nil {
		return nil, err
	}

	return o.readOutbox(s)
}

func (o *Outbox) List(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
) iter.Seq2[*entity.Outbox, error] {
	return func(yield func(*entity.Outbox, error) bool) {
		for s, err := range o.fc.List(ctx, tx, tn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(o.readOutbox(s)) {
				return
			}
		}
	}
}

func (*Outbox) readOutbox(s *span.Span) (*entity.Outbox, error) {
	defer s.Release()

	o := entity.Outbox{}
	return &o, o.Read(s)
}

func (o *Outbox) Insert(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	k := key.NewKey(oid.Bytes())
	return o.fc.Insert(ctx, tx, tn, k, e)
}

func (o *Outbox) Replace(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	return o.fc.Replace(ctx,
		tx,
		tn,
		key.NewKey(oid.Bytes()),
		e)
}
