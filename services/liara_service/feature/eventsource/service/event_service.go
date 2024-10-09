package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/service"
	"github.com/cardboardrobots/eventsource/value"
)

type EventService struct {
	eventRepository  service.EventRepository
	outboxRepository service.OutboxRepository
}

func NewEventService(
	eventRepository service.EventRepository,
	outboxRepository service.OutboxRepository,
) *EventService {
	return &EventService{
		eventRepository:  eventRepository,
		outboxRepository: outboxRepository,
	}
}

func (es *EventService) Append(
	ctx context.Context,
	e ...entity.AppendEvent,
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
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetAfterGlobalVersion(ctx, version, limit)
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
