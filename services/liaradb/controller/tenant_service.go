package controller

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type TenantService interface {
	Create(
		ctx context.Context,
		cmd service.CreateTenantCommand,
	) (value.TenantID, error)

	Delete(
		ctx context.Context,
		cmd service.DeleteTenantCommand,
	) error

	Get(
		ctx context.Context,
		tenantID value.TenantID,
	) (*entity.Tenant, error)

	List(
		ctx context.Context,
		limit int,
		offset int,
	) iter.Seq2[*entity.Tenant, error]

	Rename(
		ctx context.Context,
		cmd service.RenameTenantCommand,
	) error
}
