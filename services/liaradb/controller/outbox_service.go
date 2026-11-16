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
		pid value.PartitionID,
	) (value.OutboxID, error)

	GetOutbox(
		ctx context.Context,
		tid value.TenantID,
		oid value.OutboxID,
	) (*entity.Outbox, error)

	ListOutboxes(
		ctx context.Context,
		tid value.TenantID,
	) iter.Seq2[*entity.Outbox, error]

	UpdateOutboxPosition(
		ctx context.Context,
		tid value.TenantID,
		oid value.OutboxID,
		gv value.GlobalVersion,
		sv value.SubVersion,
	) error
}
