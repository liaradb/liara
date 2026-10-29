package controller

import (
	"context"
	"iter"

	"github.com/cardboardrobots/baseerror"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type testOutboxService struct {
	outboxes map[value.OutboxID]*entity.Outbox
	events   []*entity.Event
	requests map[value.RequestID]*entity.RequestLog
	version  value.GlobalVersion
}

var _ OutboxService = (*testOutboxService)(nil)

func (es *testOutboxService) CreateOutbox(
	ctx context.Context,
	tid value.TenantID,
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	if es.outboxes == nil {
		es.outboxes = make(map[value.OutboxID]*entity.Outbox)
	}

	id := value.NewOutboxID()
	es.outboxes[id] = entity.NewOutbox(id, partitionRange)

	return id, nil
}

func (es *testOutboxService) GetOutbox(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
) (*entity.Outbox, error) {
	o, ok := es.outboxes[outboxID]
	if !ok {
		return nil, baseerror.ErrNotFound
	}

	return o, nil
}

func (es *testOutboxService) ListOutboxes(
	ctx context.Context,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	return func(yield func(*entity.Outbox, error) bool) {
		for _, o := range es.outboxes {
			if !yield(o, nil) {
				return
			}
		}
	}
}

func (es *testOutboxService) UpdateOutboxPosition(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	o, ok := es.outboxes[outboxID]
	if !ok {
		return baseerror.ErrNotFound
	}

	o.UpdateGlobalVersion(globalVersion)
	return nil
}
