package tenant

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

type Tenant struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *Tenant {
	return &Tenant{
		fc: fixed.New(s, c, l),
	}
}

func (t *Tenant) Get(
	ctx context.Context,
	l span.Log,
	tid value.TenantID,
) (*entity.Tenant, error) {
	k := key.NewKey(tid.Bytes())
	s, err := t.fc.GetByKey(ctx, l, tablename.Tenant, k)
	if err != nil {
		return nil, err
	}

	return t.readTenant(s)
}

func (t *Tenant) List(
	ctx context.Context,
	l span.Log,
) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for s, err := range t.fc.List(ctx, l, tablename.Tenant) {
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
	l span.Log,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	k := key.NewKey(tid.Bytes())
	return t.fc.Insert(ctx,
		l,
		tablename.Tenant,
		k,
		e)
}

func (t *Tenant) Replace(
	ctx context.Context,
	l span.Log,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	return t.fc.Replace(ctx,
		l,
		tablename.Tenant,
		key.NewKey(tid.Bytes()),
		e)
}
