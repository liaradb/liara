package schema

import "context"

type Manager struct {
}

func (m *Manager) GetSchema(ctx context.Context, name string) (*Schema, error) {
	panic("unimplemented")
}
