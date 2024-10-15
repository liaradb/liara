package liara

import (
	"context"
	"iter"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
)

type (
	Service[T AggregateRoot[U], U ~string] struct {
		eventRepository EventRepository
		fromEvent       func(name string, data []byte) (any, error)
		init            func() T
	}

	AggregateRoot[U ~string] interface {
		ID() U
		Apply(any) // The method to project an Event onto the Aggregate
	}
)

func NewService[T AggregateRoot[U], U ~string, E EventData](
	eventRepository EventRepository,
	fromEvent func(name string, data []byte) (E, error),
	init func() T,
) *Service[T, U] {
	return &Service[T, U]{
		eventRepository: eventRepository,
		fromEvent:       func(name string, data []byte) (any, error) { return fromEvent(name, data) },
		init:            init,
	}
}

func (s *Service[T, U]) Append(
	ctx context.Context,
	events ...EventOptions[U],
) error {
	data := make([]AppendEvent, 0, len(events))
	for _, e := range events {
		event, err := e.toAppendEvent()
		if err != nil {
			return err
		}

		data = append(data, event)
	}

	return s.eventRepository.Append(ctx, data...)
}

func (s *Service[T, U]) GetByID(
	ctx context.Context,
	id U,
) (T, value.Version, error) {
	return s.apply(s.eventRepository.Get(ctx,
		value.AggregateID(id)))
}

func (s *Service[T, U]) GetByIDAndName(
	ctx context.Context,
	id U,
	name value.AggregateName,
) (T, value.Version, error) {
	return s.apply(s.eventRepository.GetByAggregateIDAndName(ctx,
		value.AggregateID(id),
		name))
}

func (s *Service[T, U]) apply(
	events iter.Seq2[entity.Event, error],
) (T, value.Version, error) {
	t := s.init()
	var version value.Version

	for e, err := range events {
		if err != nil {
			return t, version, err
		}

		data, err := s.fromEvent(e.Name.String(), e.Data)
		if err != nil {
			return t, version, err
		}

		if version < e.Version {
			version = e.Version
		}
		t.Apply(data)
	}

	return t, version, nil
}
