package manager

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
)

const table = "tables"

type Manager struct {
	s  *storage.Storage
	kv *keyvalue.KeyValue
	tn tablename.TableName
}

func New(kv *keyvalue.KeyValue) *Manager {
	return &Manager{
		kv: kv,
		tn: tablename.NewFromString(table),
	}
}

func (m *Manager) Get(ctx context.Context, k key.Key) (int64, error) {
	d, err := m.kv.Get(ctx, m.tn, k)
	if err != nil {
		return 0, err
	}

	b := raw.NewBufferFromSlice(d)
	var i int64
	return i, raw.ReadInt64(b, &i)
}

func (m *Manager) Insert(ctx context.Context, k key.Key, i int64) error {
	b := raw.NewBuffer(8)
	raw.WriteInt64(b, i)
	return m.kv.Set(ctx, m.tn, k, b.Bytes())
}

func (m *Manager) List(ctx context.Context) ([]int64, error) {
	result := make([]int64, 0)
	for d, err := range m.kv.List(ctx, m.tn) {
		if err != nil {
			return nil, err
		}

		b := raw.NewBufferFromSlice(d)
		var i int64
		if err := raw.ReadInt64(b, &i); err != nil {
			return nil, err
		}

		result = append(result, i)
	}
	return result, nil
}
