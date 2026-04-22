package controller

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type EventService interface {
	Append(
		ctx context.Context,
		tenantID value.TenantID,
		options service.AppendOptions,
		pid value.PartitionID,
		e ...service.AppendEvent,
	) error

	CreateOutbox(
		ctx context.Context,
		tid value.TenantID,
		partitionRange value.PartitionRange,
	) (value.OutboxID, error)

	Get(
		ctx context.Context,
		tid value.TenantID,
		partitionID value.PartitionID,
		id value.AggregateID,
	) iter.Seq2[*entity.Event, error]

	GetAfterGlobalVersion(
		ctx context.Context,
		tid value.TenantID,
		version value.GlobalVersion,
		partitionRange value.PartitionRange,
		limit value.Limit,
	) iter.Seq2[*entity.Event, error]

	GetByAggregateIDAndName(
		ctx context.Context,
		tid value.TenantID,
		partitionID value.PartitionID,
		id value.AggregateID,
		name value.AggregateName,
	) iter.Seq2[*entity.Event, error]

	GetByOutbox(
		ctx context.Context,
		tid value.TenantID,
		outboxID value.OutboxID,
		limit value.Limit,
	) iter.Seq2[*entity.Event, error]

	GetOutbox(
		ctx context.Context,
		tid value.TenantID,
		outboxID value.OutboxID,
	) (*entity.Outbox, error)

	ListOutboxes(
		ctx context.Context,
		tid value.TenantID,
	) iter.Seq2[*entity.Outbox, error]

	TestIdempotency(
		ctx context.Context,
		tid value.TenantID,
		id value.RequestID,
	) (result bool, err error)

	UpdateOutboxPosition(
		ctx context.Context,
		tid value.TenantID,
		outboxID value.OutboxID,
		globalVersion value.GlobalVersion,
	) error
}
