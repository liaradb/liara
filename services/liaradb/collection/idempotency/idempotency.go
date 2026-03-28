package idempotency

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

type Idempotency struct {
	fc *fixed.FixedCollection
}

func New(storage *storage.Storage, cursor *btree.Cursor) *Idempotency {
	return &Idempotency{
		fc: fixed.New(storage, cursor),
	}
}

func (i *Idempotency) Get(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
) (*entity.RequestLog, error) {
	k := key.NewKey(rqid.Bytes())
	data, err := i.fc.Get(ctx, tn.RequestLog(), tn.Index(0, pid), pid, k)
	if err != nil {
		return nil, err
	}

	e := &entity.RequestLog{}
	if _, ok := e.Read(data); !ok {
		return nil, io.EOF
	}

	return e, nil
}

func (i *Idempotency) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.RequestLog, error] {
	return func(yield func(*entity.RequestLog, error) bool) {
		for data, err := range i.fc.List(ctx, tn.RequestLog(), tn.Index(0, pid), pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			e := &entity.RequestLog{}
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

func (i *Idempotency) Set(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
	e *entity.RequestLog,
) error {
	v := make([]byte, entity.RequestLogSize)
	if _, ok := e.Write(v); !ok {
		return io.EOF
	}

	k := key.NewKey(rqid.Bytes())
	return i.fc.Set(ctx, tn.RequestLog(), tn.Index(0, pid), k, v)
}

func (i *Idempotency) Test(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
) (bool, error) {
	k := key.NewKey(rqid.Bytes())
	return i.fc.Test(ctx, tn.RequestLog(), tn.Index(0, pid), k)
}
