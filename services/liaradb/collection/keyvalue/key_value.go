package keyvalue

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/log"
)

type KeyValue struct {
	fc *fixed.FixedCollection
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *KeyValue {
	return &KeyValue{
		fc: fixed.New(fixed.NewRecordCollection(s, l), c),
	}
}

func (kv *KeyValue) Get(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
) ([]byte, error) {
	s, err := kv.fc.GetByKey(ctx, tx, tn, k)
	if err != nil {
		return nil, err
	}

	return s.BytesAndRelease()
}

func (kv *KeyValue) List(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for s, err := range kv.fc.List(ctx, tx, tn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(s.BytesAndRelease()) {
				return
			}
		}
	}
}

func (kv *KeyValue) Set(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	k key.Key,
	w serializer.SizedSerializer,
) error {
	return kv.fc.Insert(ctx, tx, tn, k, w)
}
