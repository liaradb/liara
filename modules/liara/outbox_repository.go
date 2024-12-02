package liara

import (
	"context"
)

type OutboxRepository interface {
	CreateOutbox(context.Context, OutboxID, []PartitionID) (OutboxID, error)
	GetOutbox(context.Context, OutboxID) (GlobalVersion, error)
	UpdateOutboxPosition(context.Context, OutboxID, GlobalVersion) error
}
