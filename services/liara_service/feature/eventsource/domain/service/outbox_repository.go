package service

import (
	"context"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type OutboxRepository interface {
	GetOrCreateOutbox(context.Context, value.OutboxID) (value.GlobalVersion, error)
	UpdateOutboxPosition(context.Context, value.OutboxID, value.GlobalVersion) error
}
