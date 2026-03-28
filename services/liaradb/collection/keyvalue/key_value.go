package keyvalue

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage"
)

type KeyValue struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor) *KeyValue {
	return &KeyValue{
		fc: fixed.New(s, c),
	}
}

func (kv *KeyValue) Get(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	k key.Key,
) ([]byte, error) {
	return kv.fc.Get(ctx, tn.KeyValue(pid), tn.Index(0, pid), pid, k)
}

func (kv *KeyValue) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[[]byte, error] {
	return kv.fc.List(ctx, tn.KeyValue(pid), tn.Index(0, value.NewPartitionID(0)), pid)
}

func (kv *KeyValue) Set(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	k key.Key,
	v []byte,
) error {
	return kv.fc.Set(ctx, tn.KeyValue(pid), tn.Index(0, pid), k, v)
}
