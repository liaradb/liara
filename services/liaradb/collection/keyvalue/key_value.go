package keyvalue

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/serializer"
)

type KeyValue struct {
	ic *indexed.IndexedCollection
}

func New(ic *indexed.IndexedCollection) *KeyValue {
	return &KeyValue{ic: ic}
}

func (kv *KeyValue) Get(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	k key.Key,
) ([]byte, error) {
	s, err := kv.ic.GetByKey(ctx, tx, tn, k)
	if err != nil {
		return nil, err
	}

	return s.BytesAndRelease()
}

func (kv *KeyValue) List(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for s, err := range kv.ic.List(ctx, tx, tn) {
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
	tx indexed.Transaction,
	tn tablename.TableName,
	k key.Key,
	w serializer.SizedSerializer,
) error {
	_, err := kv.ic.Insert(ctx, tx, tn, k, w)
	return err
}
