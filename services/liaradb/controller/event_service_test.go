package controller

import (
	"context"
	"iter"
	"slices"
	"time"

	"github.com/cardboardrobots/baseerror"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/query"
	"github.com/liaradb/liaradb/domain/value"
)

type testEventService struct {
	outboxes map[value.OutboxID]*entity.Outbox
	events   []*entity.Event
	requests map[value.RequestID]*entity.RequestLog
	version  value.GlobalVersion
}

var _ EventService = (*testEventService)(nil)

func (es *testEventService) Append(
	ctx context.Context,
	cmd command.AppendEvent,
) error {
	if cmd.RequestID != nil {
		if _, ok := es.requests[*cmd.RequestID]; ok {
			return baseerror.ErrAlreadyExists
		}
	}

	es.version = value.NewGlobalVersion(es.version.Value() + 1)
	var sub value.SubVersion
	for _, event := range cmd.Events {
		e, err := entity.NewEvent(event, entity.NewMetadata(cmd), es.version, sub)
		if err != nil {
			return err
		}

		es.events = append(es.events, e)
		sub = sub.Increment()
	}
	slices.SortFunc(es.events, func(a, b *entity.Event) int {
		return int(a.Version().Value() - b.Version().Value())
	})

	if cmd.RequestID != nil {
		if es.requests == nil {
			es.requests = make(map[value.RequestID]*entity.RequestLog)
		}
		es.requests[*cmd.RequestID] = entity.NewRequestLog(*cmd.RequestID, value.NewTime(time.Now()))
	}

	return nil
}

func (es *testEventService) Get(
	ctx context.Context,
	tid value.TenantID,
	partitionID value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for _, event := range es.events {
			if event.PartitionID() == partitionID &&
				event.AggregateID() == id {
				if !yield(event, nil) {
					return
				}
			}
		}
	}
}

func (es *testEventService) GetAfterGlobalVersion(
	ctx context.Context,
	tid value.TenantID,
	version value.GlobalVersion,
	sub value.SubVersion,
	pid value.PartitionID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	panic("unimplemented")
}

func (es *testEventService) GetByAggregateIDAndName(
	ctx context.Context,
	tid value.TenantID,
	pid value.PartitionID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for _, event := range es.events {
			if event.PartitionID() == pid &&
				event.AggregateID() == id &&
				event.AggregateName() == name {
				if !yield(event, nil) {
					return
				}
			}
		}
	}
}

func (es *testEventService) GetByOutbox(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	panic("unimplemented")
}

func (es *testEventService) TestIdempotency(
	ctx context.Context,
	qry query.TestIdempotency,
) (result bool, err error) {
	_, ok := es.requests[qry.RequestID]
	return !ok, nil
}
