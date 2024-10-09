package service

import (
	"context"

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
		// The method to project an Event onto the Aggregate
		Apply(any)
	}
)

func NewService[T AggregateRoot[U], U ~string, E entity.EventData](
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
	correlationID value.CorrelationID,
	userID value.UserID,
	id U,
	version value.Version,
	events ...entity.EventData,
) error {
	options := entity.EventOptions{
		AggregateID:   value.AggregateID(id),
		Version:       version,
		CorrelationID: correlationID,
		UserID:        userID,
	}

	data := make([]entity.AppendEvent, 0, len(events))
	for _, item := range events {
		event, err := entity.NewEvent(options, item)
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
	t := s.init()
	var version value.Version

	for e, err := range s.eventRepository.Get(ctx, value.AggregateID(id)) {
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

func (s *Service[T, U]) GetByIDAndName(
	ctx context.Context,
	id U,
	name value.AggregateName,
) (T, value.Version, error) {
	t := s.init()
	var version value.Version

	for e, err := range s.eventRepository.GetByAggregateIDAndName(ctx, value.AggregateID(id), name) {
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

// func (s *Service[T, U]) aggregateCallback(callback func(Event, any) error) func(em Event) error {
// 	return func(em Event) error {
// 		data, err := s.fromEvent(em.Name.String(), em.Data)
// 		if err != nil {
// 			return err
// 		}

// 		return callback(em, data)
// 	}
// }
