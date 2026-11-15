package outbox

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap/span"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type Outbox struct {
	ic *indexed.IndexedCollection
}

func New(ic *indexed.IndexedCollection) *Outbox {
	return &Outbox{ic: ic}
}

func (o *Outbox) Get(
	ctx context.Context,
	tx indexed.Transaction,
	tid value.TenantID,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	tn := tablename.NewTenantFolder(tid).Outbox()
	k := key.NewKey(oid.Bytes())
	s, err := o.ic.GetByKey(ctx, tx, tn, k)
	if err != nil {
		return nil, err
	}

	return o.readOutbox(s)
}

func (o *Outbox) List(
	ctx context.Context,
	tx indexed.Transaction,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	tn := tablename.NewTenantFolder(tid).Outbox()
	return func(yield func(*entity.Outbox, error) bool) {
		for s, err := range o.ic.List(ctx, tx, tn) {
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
	tx indexed.Transaction,
	tid value.TenantID,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	tn := tablename.NewTenantFolder(tid).Outbox()
	k := key.NewKey(oid.Bytes())
	_, err := o.ic.Insert(ctx, tx, tn, k, e)
	return err
}

func (o *Outbox) Replace(
	ctx context.Context,
	tx indexed.Transaction,
	tid value.TenantID,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	tn := tablename.NewTenantFolder(tid).Outbox()
	return o.ic.Replace(ctx,
		tx,
		tn,
		key.NewKey(oid.Bytes()),
		e)
}
