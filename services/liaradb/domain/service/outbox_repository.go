package service

import (
	"context"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type OutboxRepository interface {
	CreateTable(context.Context, value.TenantID) error
	CreateOutbox(context.Context, value.TenantID, *entity.Outbox) error
	DropTable(context.Context, value.TenantID) error
	GetOutbox(context.Context, value.TenantID, value.OutboxID) (*entity.Outbox, error)
	UpdateOutboxPosition(context.Context, value.TenantID, value.OutboxID, value.GlobalVersion) error
}
