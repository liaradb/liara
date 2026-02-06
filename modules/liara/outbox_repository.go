package liara

import (
	"context"
)

type OutboxRepository interface {
	CreateOutbox(context.Context, TenantID, PartitionID, PartitionID) (OutboxID, error)
	GetOutbox(context.Context, TenantID, OutboxID) (GlobalVersion, error)
	UpdateOutboxPosition(context.Context, TenantID, OutboxID, GlobalVersion) error
}
