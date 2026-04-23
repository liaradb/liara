package controller

import (
	"context"
	"iter"
	"slices"

	"github.com/cardboardrobots/baseerror"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type testEventService struct {
	outboxes map[value.OutboxID]*entity.Outbox
	events   []*entity.Event
	version  value.GlobalVersion
}

func (es *testEventService) Append(ctx context.Context, tenantID value.TenantID, options service.AppendOptions, pid value.PartitionID, e ...service.AppendEvent) error {
	for _, event := range e {
		id, err := value.NewEventIDFromString(event.ID)
		if err != nil {
			return err
		}

		es.version = value.NewGlobalVersion(es.version.Value() + 1)
		es.events = append(es.events, &entity.Event{
			GlobalVersion: es.version,
			PartitionID:   pid,
			AggregateID:   event.AggregateID,
			ID:            id,
			Version:       event.Version,
		})
	}
	slices.SortFunc(es.events, func(a, b *entity.Event) int {
		return int(a.Version.Value() - b.Version.Value())
	})
	return nil
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
	return func(yield func(*entity.Event, error) bool) {
		for _, event := range es.events {
			if event.PartitionID == partitionID && event.AggregateID == id {
				if !yield(event, nil) {
					return
				}
			}
		}
	}
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
	return func(yield func(*entity.Outbox, error) bool) {
		for _, o := range es.outboxes {
			if !yield(o, nil) {
				return
			}
		}
	}
}

func (es *testEventService) TestIdempotency(ctx context.Context, tid value.TenantID, id value.RequestID) (result bool, err error) {
	panic("unimplemented")
}

func (es *testEventService) UpdateOutboxPosition(ctx context.Context, tid value.TenantID, outboxID value.OutboxID, globalVersion value.GlobalVersion) error {
	o, ok := es.outboxes[outboxID]
	if !ok {
		return baseerror.ErrNotFound
	}

	o.UpdateGlobalVersion(globalVersion)
	return nil
}

var _ EventService = (*testEventService)(nil)
