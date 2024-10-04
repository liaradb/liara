package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/eventsource"
)

type EventService struct {
	eventSource     eventsource.EventSource
	eventRepository eventsource.EventRepository
}

func NewEventService(
	eventSource eventsource.EventSource,
	eventRepository eventsource.EventRepository,
) *EventService {
	return &EventService{
		eventSource:     eventSource,
		eventRepository: eventRepository,
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
