package schema

import (
	"context"

	"github.com/liaradb/liaradb/storage"
)

type Manager struct {
	s *storage.Storage
}

func NewManager(s *storage.Storage) *Manager {
	return &Manager{
		s: s,
	}
}

func (m *Manager) GetSchema(ctx context.Context, name string) (*Schema, error) {
	panic("unimplemented")
}
