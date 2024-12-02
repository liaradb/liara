package service

import (
	"context"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type OutboxRepository interface {
	CreateOutbox(context.Context, *entity.Outbox) error
	GetOutbox(context.Context, value.OutboxID) (*entity.Outbox, error)
	UpdateOutboxPosition(context.Context, value.OutboxID, value.GlobalVersion) error
}
