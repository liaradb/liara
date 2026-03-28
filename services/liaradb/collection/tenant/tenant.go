package tenant

import (
	"context"
	"io"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage"
)

type Tenant struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor) *Tenant {
	return &Tenant{
		fc: fixed.New(s, c),
	}
}

func (t *Tenant) Get(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	tid value.TenantID,
) (*entity.Tenant, error) {
	k := key.NewKey(tid.Bytes())
	data, err := t.fc.Get(ctx, tn.RequestLog(), tn.Index(0, pid), k)
	if err != nil {
		return nil, err
	}

	e := &entity.Tenant{}
	if _, ok := e.Read(data); !ok {
		return nil, io.EOF
	}

	return e, nil
}

func (t *Tenant) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for data, err := range t.fc.List(ctx, tn.RequestLog(), tn.Index(0, pid), pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			e := &entity.Tenant{}
			if _, ok := e.Read(data); !ok {
				yield(nil, io.EOF)
				return
			}

			if !yield(e, nil) {
				return
			}
		}
	}
}

func (t *Tenant) Set(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	v := make([]byte, entity.TenantSize)
	if _, ok := e.Write(v); !ok {
		return io.EOF
	}

	k := key.NewKey(tid.Bytes())
	return t.fc.Set(ctx, tn.RequestLog(), tn.Index(0, pid), k, v)
}

func (t *Tenant) Replace(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	v := make([]byte, entity.TenantSize)
	if _, ok := e.Write(v); !ok {
		return btree.ErrNoUpdate
	}

	return t.fc.Replace(ctx,
		tn.Outbox(pid),
		tn.Index(0, value.NewPartitionID(0)),
		pid,
		key.NewKey(tid.Bytes()),
		v)
}
