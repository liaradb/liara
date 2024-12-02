package liara

import (
	"context"
)

type OutboxRepository interface {
	GetOrCreateOutbox(context.Context, OutboxID, PartitionID) (GlobalVersion, error)
	UpdateOutboxPosition(context.Context, OutboxID, GlobalVersion) error
}
