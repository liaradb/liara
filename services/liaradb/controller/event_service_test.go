package controller

import (
	"context"
	"iter"

	"github.com/cardboardrobots/baseerror"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type testEventService struct {
	outboxes map[value.OutboxID]*entity.Outbox
}

func (es *testEventService) Append(ctx context.Context, tenantID value.TenantID, options service.AppendOptions, pid value.PartitionID, e ...service.AppendEvent) error {
	panic("unimplemented")
}

func (es *testEventService) CreateOutbox(ctx context.Context, tid value.TenantID, partitionRange value.PartitionRange) (value.OutboxID, error) {
	if es.outboxes == nil {
		es.outboxes = make(map[value.OutboxID]*entity.Outbox)
	}

	id := value.NewOutboxID()
	es.outboxes[id] = entity.NewOutbox(id, partitionRange)

	return id, nil
}

func (es *testEventService) Get(ctx context.Context, tid value.TenantID, partitionID value.PartitionID, id value.AggregateID) iter.Seq2[*entity.Event, error] {
	panic("unimplemented")
}

func (es *testEventService) GetAfterGlobalVersion(ctx context.Context, tid value.TenantID, version value.GlobalVersion, partitionRange value.PartitionRange, limit value.Limit) iter.Seq2[*entity.Event, error] {
	panic("unimplemented")
}

func (es *testEventService) GetByAggregateIDAndName(ctx context.Context, tid value.TenantID, partitionID value.PartitionID, id value.AggregateID, name value.AggregateName) iter.Seq2[*entity.Event, error] {
	panic("unimplemented")
}

func (es *testEventService) GetByOutbox(ctx context.Context, tid value.TenantID, outboxID value.OutboxID, limit value.Limit) iter.Seq2[*entity.Event, error] {
	panic("unimplemented")
}

func (es *testEventService) GetOutbox(ctx context.Context, tid value.TenantID, outboxID value.OutboxID) (*entity.Outbox, error) {
	o, ok := es.outboxes[outboxID]
	if !ok {
		return nil, baseerror.ErrNotFound
	}

	return o, nil
}

func (es *testEventService) ListOutboxes(ctx context.Context, tid value.TenantID) iter.Seq2[*entity.Outbox, error] {
	panic("unimplemented")
}

func (es *testEventService) TestIdempotency(ctx context.Context, tid value.TenantID, id value.RequestID) (result bool, err error) {
	panic("unimplemented")
}

func (es *testEventService) UpdateOutboxPosition(ctx context.Context, tid value.TenantID, outboxID value.OutboxID, globalVersion value.GlobalVersion) error {
	panic("unimplemented")
}

var _ EventService = (*testEventService)(nil)
