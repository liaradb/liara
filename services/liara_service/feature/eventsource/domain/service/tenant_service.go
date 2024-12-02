package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
)

type TenantService struct {
	tenantRepository TenantRepository
}

type TenantRepository interface {
	Insert(context.Context, *entity.Tenant) error
	List(context.Context, int, int) iter.Seq2[*entity.Tenant, error]
}

func NewTenantService(
	repository TenantRepository,
) *TenantService {
	return &TenantService{
		tenantRepository: repository,
	}
}

func (ts *TenantService) Insert(ctx context.Context, tenant *entity.Tenant) error {
	return ts.tenantRepository.Insert(ctx, tenant)
}

func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return ts.tenantRepository.List(ctx, limit, offset)
}
