package service

import (
	"context"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type OutboxRepository interface {
	CreateTable(context.Context, value.TenantID) error
	CreateOutbox(context.Context, value.TenantID, *entity.Outbox) error
	GetOutbox(context.Context, value.TenantID, value.OutboxID) (*entity.Outbox, error)
	UpdateOutboxPosition(context.Context, value.TenantID, value.OutboxID, value.GlobalVersion) error
}
