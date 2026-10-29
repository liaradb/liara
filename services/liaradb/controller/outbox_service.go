package controller

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type OutboxService interface {
	CreateOutbox(
		ctx context.Context,
		tid value.TenantID,
		partitionRange value.PartitionRange,
	) (value.OutboxID, error)

	GetOutbox(
		ctx context.Context,
		tid value.TenantID,
		outboxID value.OutboxID,
	) (*entity.Outbox, error)

	ListOutboxes(
		ctx context.Context,
		tid value.TenantID,
	) iter.Seq2[*entity.Outbox, error]

	UpdateOutboxPosition(
		ctx context.Context,
		tid value.TenantID,
		outboxID value.OutboxID,
		globalVersion value.GlobalVersion,
	) error
}
