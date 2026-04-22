package controller

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type testTenantService struct {
}

func (ts *testTenantService) Create(ctx context.Context, cmd service.CreateTenantCommand) (value.TenantID, error) {
	panic("unimplemented")
}

func (ts *testTenantService) Delete(ctx context.Context, cmd service.DeleteTenantCommand) error {
	panic("unimplemented")
}

func (ts *testTenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	panic("unimplemented")
}

func (ts *testTenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	panic("unimplemented")
}

func (ts *testTenantService) Rename(ctx context.Context, cmd service.RenameTenantCommand) error {
	panic("unimplemented")
}

var _ TenantService = (*testTenantService)(nil)
