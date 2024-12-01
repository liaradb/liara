package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type EventService struct {
	eventRepository  EventRepository
	outboxRepository OutboxRepository
}

func NewEventService(
	eventRepository EventRepository,
	outboxRepository OutboxRepository,
) *EventService {
	return &EventService{
		eventRepository:  eventRepository,
		outboxRepository: outboxRepository,
	}
}

func (es *EventService) Append(
	ctx context.Context,
	e ...AppendEvent,
) error {
	return es.eventRepository.Append(ctx, e...)
}

func (es *EventService) Get(
	ctx context.Context,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.Get(ctx, id)
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetByAggregateIDAndName(ctx, id, name)
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	version value.GlobalVersion,
	partitionID value.PartitionID,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetAfterGlobalVersion(ctx, version, partitionID, limit)
}

func (es *EventService) GetOrCreateOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) (value.GlobalVersion, error) {
	return es.outboxRepository.GetOrCreateOutbox(ctx, outboxID)
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	return es.outboxRepository.UpdateOutboxPosition(ctx, outboxID, globalVersion)
}
