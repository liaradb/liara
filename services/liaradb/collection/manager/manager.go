package manager

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
)

type Manager struct {
	s  *storage.Storage
	kv *keyvalue.KeyValue
	tn tablename.TableName
}

func New(kv *keyvalue.KeyValue) *Manager {
	return &Manager{
		kv: kv,
		tn: tablename.Manager,
	}
}

func (m *Manager) Get(ctx context.Context, tx fixed.Transaction, k key.Key) (int64, error) {
	d, err := m.kv.Get(ctx, tx, m.tn, k)
	if err != nil {
		return 0, err
	}

	b := buffer.NewFromSlice(d)
	var i int64
	return i, raw.ReadInt64(b, &i)
}

func (m *Manager) Insert(ctx context.Context, tx fixed.Transaction, k key.Key, i int64) error {
	return m.kv.Set(ctx, tx, m.tn, k, base.NewInt64(i))
}

func (m *Manager) List(ctx context.Context, tx fixed.Transaction) ([]int64, error) {
	result := make([]int64, 0)
	for d, err := range m.kv.List(ctx, tx, m.tn) {
		if err != nil {
			return nil, err
		}

		b := buffer.NewFromSlice(d)
		var i int64
		if err := raw.ReadInt64(b, &i); err != nil {
			return nil, err
		}

		result = append(result, i)
	}
	return result, nil
}
