package service

import (
	"context"

	"github.com/cardboardrobots/eventsource/value"
)

type OutboxRepository interface {
	GetOrCreateOutbox(context.Context, value.OutboxID) (value.GlobalVersion, error)
	UpdateOutboxPosition(context.Context, value.OutboxID, value.GlobalVersion) error
}
