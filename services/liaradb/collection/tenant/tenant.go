package tenant

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
	"github.com/liaradb/liaradb/transaction/log"
)

type Tenant struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *Tenant {
	return &Tenant{
		fc: fixed.New(record.NewRecordCollection(s, l), c),
	}
}

func (t *Tenant) Get(
	ctx context.Context,
	tx span.Transaction,
	tid value.TenantID,
) (*entity.Tenant, error) {
	k := key.NewKey(tid.Bytes())
	s, err := t.fc.GetByKey(ctx, tx, tablename.Tenant, k)
	if err != nil {
		return nil, err
	}

	return t.readTenant(s)
}

func (t *Tenant) List(
	ctx context.Context,
	tx span.Transaction,
) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for s, err := range t.fc.List(ctx, tx, tablename.Tenant) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(t.readTenant(s)) {
				return
			}
		}
	}
}

func (*Tenant) readTenant(s *span.Span) (*entity.Tenant, error) {
	defer s.Release()

	t := entity.Tenant{}
	return &t, t.Read(s)
}

func (t *Tenant) Insert(
	ctx context.Context,
	tx span.Transaction,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	k := key.NewKey(tid.Bytes())
	_, err := t.fc.Insert(ctx,
		tx,
		tablename.Tenant,
		k,
		e)
	return err
}

func (t *Tenant) Replace(
	ctx context.Context,
	tx span.Transaction,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	return t.fc.Replace(ctx,
		tx,
		tablename.Tenant,
		key.NewKey(tid.Bytes()),
		e)
}
