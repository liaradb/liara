package controller

import (
	"context"
	"iter"

	"github.com/cardboardrobots/baseerror"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type testTenantService struct {
	tenants map[value.TenantID]*entity.Tenant
}

var _ TenantService = (*testTenantService)(nil)

func (ts *testTenantService) Create(ctx context.Context, cmd service.CreateTenantCommand) (value.TenantID, error) {
	id := value.NewTenantID()
	if ts.tenants == nil {
		ts.tenants = make(map[value.TenantID]*entity.Tenant)
	}
	ts.tenants[id] = entity.NewTenant(id, cmd.TenantName)
	return id, nil
}

func (ts *testTenantService) Delete(ctx context.Context, cmd service.DeleteTenantCommand) error {
	_, ok := ts.tenants[cmd.TenantID]
	if !ok {
		return baseerror.ErrNotFound
	}

	delete(ts.tenants, cmd.TenantID)
	return nil
}

func (ts *testTenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	t, ok := ts.tenants[tenantID]
	if !ok {
		return nil, baseerror.ErrNotFound
	}

	return t, nil
}

func (ts *testTenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for _, t := range ts.tenants {
			if !yield(t, nil) {
				return
			}
		}
	}
}

func (ts *testTenantService) Rename(ctx context.Context, cmd service.RenameTenantCommand) error {
	t, ok := ts.tenants[cmd.TenantID]
	if !ok {
		return baseerror.ErrNotFound
	}

	return t.Rename(cmd.TenantName)
}
