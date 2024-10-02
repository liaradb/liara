package eventsource

import (
	"context"
	"iter"
)

type (
	Service[T AggregateRoot[U], U ~string] struct {
		eventSource EventSource
		fromEvent   func(name string, data []byte) (any, error)
		init        func() T
	}

	AggregateRoot[U ~string] interface {
		ID() U
		// The method to project an Event onto the Aggregate
		Apply(any)
	}

	EventSource interface {
		Get(ctx context.Context, id AggregateID) iter.Seq2[Event, error]
		GetByAggregateIDAndName(ctx context.Context, id AggregateID, name AggregateName) iter.Seq2[Event, error]
		Append(ctx context.Context, e ...Event) error
	}
)

func NewService[T AggregateRoot[U], U ~string, E EventData](
	eventSource EventSource,
	fromEvent func(name string, data []byte) (E, error),
	init func() T,
) *Service[T, U] {
	return &Service[T, U]{
		eventSource: eventSource,
		fromEvent:   func(name string, data []byte) (any, error) { return fromEvent(name, data) },
		init:        init,
	}
}

func (s *Service[T, U]) Append(
	ctx context.Context,
	correlationID CorrelationID,
	userID UserID,
	id U,
	version Version,
	events ...EventData,
) error {
	options := EventOptions{
		AggregateID:   AggregateID(id),
		Version:       version,
		CorrelationID: correlationID,
		UserID:        userID,
	}

	data := make([]Event, 0, len(events))
	for _, item := range events {
		event, err := newEvent(options, item)
		if err != nil {
			return err
		}

		data = append(data, event)
	}

	return s.eventSource.Append(ctx, data...)
}

func (s *Service[T, U]) GetByID(
	ctx context.Context,
	id U,
) (T, Version, error) {
	t := s.init()
	var version Version

	rows := s.eventSource.Get(ctx, AggregateID(id))

	for e, err := range rows {
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
	name AggregateName,
) (T, Version, error) {
	t := s.init()
	var version Version

	rows := s.eventSource.GetByAggregateIDAndName(ctx, AggregateID(id), name)

	for e, err := range rows {
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
