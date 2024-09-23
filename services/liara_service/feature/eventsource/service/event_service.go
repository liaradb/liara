package service

import (
	"context"
	"errors"
	"iter"

	"github.com/cardboardrobots/eventsource"
)

type EventService struct {
	eventSource eventsource.EventSource
}

func NewEventService(
	eventSource eventsource.EventSource,
) *EventService {
	return &EventService{
		eventSource: eventSource,
	}
}

func (es *EventService) Get(ctx context.Context) iter.Seq2[eventsource.Event, error] {
	return func(yield func(eventsource.Event, error) bool) {
		yield(eventsource.Event{}, errors.New("unimplemented"))
	}
}
