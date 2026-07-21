package controller

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type TenantService interface {
	Create(
		ctx context.Context,
		cmd command.CreateTenant,
	) (value.TenantID, error)

	Delete(
		ctx context.Context,
		cmd command.DeleteTenant,
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
		cmd command.RenameTenant,
	) error
}
