package tenant

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/infrastructure/tenantfolder"
	"github.com/liaradb/liaradb/domain/value"
)

type Tenant struct {
	ic *indexed.IndexedCollection
}

func New(ic *indexed.IndexedCollection) *Tenant {
	return &Tenant{ic: ic}
}

func (t *Tenant) Get(
	ctx context.Context,
	tx indexed.Transaction,
	tid value.TenantID,
) (*entity.Tenant, error) {
	k := key.NewKey(tid.Bytes())
	r, err := t.ic.GetByKey(ctx, tx, tenantfolder.Tenant, k)
	if err != nil {
		return nil, err
	}

	return t.readTenant(r)
}

func (t *Tenant) List(
	ctx context.Context,
	tx indexed.Transaction,
) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for r, err := range t.ic.List(ctx, tx, tenantfolder.Tenant) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(t.readTenant(r)) {
				return
			}
		}
	}
}

func (*Tenant) readTenant(s *record.Record) (*entity.Tenant, error) {
	defer s.Release()

	t := entity.Tenant{}
	return &t, t.Read(s)
}

func (t *Tenant) Insert(
	ctx context.Context,
	tx indexed.Transaction,
	e *entity.Tenant,
) error {
	k := key.NewKey(e.ID().Bytes())
	_, err := t.ic.Insert(ctx,
		tx,
		tenantfolder.Tenant,
		k,
		e)
	return err
}

func (t *Tenant) Replace(
	ctx context.Context,
	tx indexed.Transaction,
	e *entity.Tenant,
) error {
	return t.ic.Replace(ctx,
		tx,
		tenantfolder.Tenant,
		key.NewKey(e.ID().Bytes()),
		e)
}
