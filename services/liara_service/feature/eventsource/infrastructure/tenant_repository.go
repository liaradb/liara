package infrastructure

import (
	"context"
	"iter"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
)

type TenantRepository struct {
}

var _ service.TenantRepository = (*TenantRepository)(nil)

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{}
}

func (t *TenantRepository) Insert(context.Context, *entity.Tenant) error {
	return nil
}

func (t *TenantRepository) List(context.Context, int, int) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {}
}
