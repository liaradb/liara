package transaction

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/util/set"
)

type keyValueTransaction struct {
	values []valueItem
	keys   set.Set[key.Key]
	kv     *keyvalue.KeyValue
}

type valueItem struct {
	r    *entity.Row
	data []byte
}

func newKeyValueTransaction(kv *keyvalue.KeyValue) keyValueTransaction {
	return keyValueTransaction{
		kv:   kv,
		keys: set.Set[key.Key]{},
	}
}

func (kv *keyValueTransaction) setValue(
	ctx context.Context,
	tn tablename.TableName,
	r *entity.Row,
	data []byte,
) error {
	// k := key.NewKey2(r.ID().Bytes(), r.Version().Value())
	// // Verify this AggregateID and Version is unique in this transaction
	// if t.keys.Includes(k) {
	// 	return btree.ErrExists
	// }

	// t.keys.Add(k)

	// // Verify this AggregateID and Version is unique in Index
	// if err := t.keyValue.CanAppend(ctx, tn, r.PartitionID(), k); err != nil {
	// 	return err
	// }

	kv.values = append(kv.values, valueItem{
		r:    r,
		data: data,
	})
	return nil
}
