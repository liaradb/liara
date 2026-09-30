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
	tn tablename.TableName,
	pid value.PartitionID,
	tid value.TenantID,
) (*entity.Tenant, error) {
	k := key.NewKey(tid.Bytes())
	data, err := t.fc.Get(ctx, tn.Tenant(), tn.Index(0, pid), k)
	if err != nil {
		return nil, err
	}

	e := &entity.Tenant{}
	_ = e.Read(data)
	return e, nil
}

func (t *Tenant) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for data, err := range t.fc.List(ctx, tn.Tenant(), tn.Index(0, pid), pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			e := &entity.Tenant{}
			_ = e.Read(data)
			if !yield(e, nil) {
				return
			}
		}
	}
}

func (t *Tenant) Set(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	pid value.PartitionID,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	v := make([]byte, entity.TenantSize)
	_ = e.Write(v)
	k := key.NewKey(tid.Bytes())
	return t.fc.Insert(ctx, l, tn.Tenant(), tn.Index(0, pid), k, v)
}

func (t *Tenant) Replace(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	pid value.PartitionID,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	v := make([]byte, entity.TenantSize)
	_ = e.Write(v)
	return t.fc.Replace(ctx,
		l,
		tn.Tenant(),
		tn.Index(0, pid),
		pid,
		key.NewKey(tid.Bytes()),
		v)
}
