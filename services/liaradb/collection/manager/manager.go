package manager

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
)

type Manager struct {
	s  *storage.Storage
	kv *keyvalue.KeyValue
}

func New(s *storage.Storage) *Manager {
	return &Manager{
		s:  s,
		kv: keyvalue.New(s),
	}
}

func (m *Manager) Insert(ctx context.Context, k value.Key, i int64) error {
	b := raw.NewBuffer(8)
	raw.WriteInt64(b, i)
	return m.kv.Set(ctx, tablename.New("tables"), k, b.Bytes())
}

func (m *Manager) List(ctx context.Context, k value.Key) (int64, error) {
	d, err := m.kv.Get(ctx, tablename.New("tables"), k)
	if err != nil {
		return 0, err
	}

	b := raw.NewBufferFromSlice(d)
	var i int64
	return i, raw.ReadInt64(b, &i)
}
