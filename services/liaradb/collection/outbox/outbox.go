package outbox

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

type Outbox struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor) *Outbox {
	return &Outbox{
		fc: fixed.New(s, c),
	}
}

func (o *Outbox) Get(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	k := key.NewKey(oid.Bytes())
	data, err := o.fc.Get(ctx, tn.RequestLog(), tn.Index(0, pid), k)
	if err != nil {
		return nil, err
	}

	e := &entity.Outbox{}
	if _, ok := e.Read(data); !ok {
		return nil, io.EOF
	}

	return e, nil
}

func (o *Outbox) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Outbox, error] {
	return func(yield func(*entity.Outbox, error) bool) {
		for data, err := range o.fc.List(ctx, tn.RequestLog(), tn.Index(0, pid), pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			e := &entity.Outbox{}
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

func (o *Outbox) Set(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	v := make([]byte, entity.OutboxSize)
	if _, ok := e.Write(v); !ok {
		return io.EOF
	}

	k := key.NewKey(oid.Bytes())
	return o.fc.Set(ctx, tn.RequestLog(), tn.Index(0, pid), k, v)
}

func (o *Outbox) Replace(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	v := make([]byte, entity.OutboxSize)
	if _, ok := e.Write(v); !ok {
		return btree.ErrNoUpdate
	}

	return o.fc.Replace(ctx,
		tn.Outbox(pid),
		tn.Index(0, value.NewPartitionID(0)),
		pid,
		key.NewKey(oid.Bytes()),
		v)
}
