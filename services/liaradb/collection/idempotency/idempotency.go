package idempotency

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

type Idempotency struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *Idempotency {
	return &Idempotency{
		fc: fixed.New(s, c, l),
	}
}

func (i *Idempotency) Get(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
) (*entity.RequestLog, error) {
	k := key.NewKey(rqid.Bytes())
	s, err := i.fc.GetByKey(ctx, tn.RequestLog(), tn.Index(0, pid), k)
	if err != nil {
		return nil, err
	}

	return i.readRequestLog(s)
}

func (i *Idempotency) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.RequestLog, error] {
	return func(yield func(*entity.RequestLog, error) bool) {
		for s, err := range i.fc.List(ctx, tn.RequestLog(), tn.Index(0, pid), pid) {
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

func (i *Idempotency) Set(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
	e *entity.RequestLog,
) error {
	v := make([]byte, entity.RequestLogSize)
	_ = e.Write(v)

	k := key.NewKey(rqid.Bytes())
	return i.fc.Insert(ctx, l, tn.RequestLog(), tn.Index(0, pid), k, v)
}

func (i *Idempotency) Test(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
) (bool, error) {
	k := key.NewKey(rqid.Bytes())
	return i.fc.Test(ctx, tn.Index(0, pid), k)
}
