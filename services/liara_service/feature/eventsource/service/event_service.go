package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/eventsource"
)

type EventService struct {
	eventSource      eventsource.EventSource
	eventRepository  eventsource.EventRepository
	outboxRepository eventsource.OutboxRepository
}

func NewEventService(
	eventSource eventsource.EventSource,
	eventRepository eventsource.EventRepository,
	outboxRepository eventsource.OutboxRepository,
) *EventService {
	return &EventService{
		eventSource:      eventSource,
		eventRepository:  eventRepository,
		outboxRepository: outboxRepository,
	}
}

func (es *EventService) Append(
	ctx context.Context,
	e ...eventsource.Event,
) error {
	return es.eventSource.Append(ctx, e...)
}

func (es *EventService) Get(
	ctx context.Context,
	id eventsource.AggregateID,
) iter.Seq2[eventsource.Event, error] {
	return es.eventSource.Get(ctx, id)
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	id eventsource.AggregateID,
	name eventsource.AggregateName,
) iter.Seq2[eventsource.Event, error] {
	return es.eventSource.GetByAggregateIDAndName(ctx, id, name)
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	version eventsource.GlobalVersion,
	limit eventsource.Limit,
) iter.Seq2[eventsource.Event, error] {
	return es.eventRepository.GetAfterGlobalVersion(ctx, version, limit)
}

func (es *EventService) GetOrCreateOutbox(
	ctx context.Context,
	outboxID eventsource.OutboxID,
) (eventsource.GlobalVersion, error) {
	return es.outboxRepository.GetOrCreateOutbox(ctx, outboxID)
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	outboxID eventsource.OutboxID,
	globalVersion eventsource.GlobalVersion,
) error {
	return es.outboxRepository.UpdateOutboxPosition(ctx, outboxID, globalVersion)
}
