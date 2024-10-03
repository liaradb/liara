package eventsource

import (
	"context"
	"iter"
)

type MockEventSource struct {
	events   []Event
	versions map[AggregateID]Version
}

var _ EventSource = &MockEventSource{}
var _ EventRepository = &MockEventSource{}

func (mes *MockEventSource) Get(
	ctx context.Context,
	id AggregateID,
) iter.Seq2[Event, error] {
	return func(yield func(Event, error) bool) {
		for _, e := range mes.events {
			if e.AggregateID != id {
				continue
			}

			if !yield(e, nil) {
				return
			}
		}
	}
}

func (mes *MockEventSource) GetAfterGlobalVersion(
	ctx context.Context,
	globalVersion GlobalVersion,
	limit Limit,
) iter.Seq2[Event, error] {
	panic("unimplemented")
}

func (mes *MockEventSource) GetByAggregateIDAndName(
	ctx context.Context,
	id AggregateID,
	name AggregateName,
) iter.Seq2[Event, error] {
	return func(yield func(Event, error) bool) {
		for _, e := range mes.events {
			if e.AggregateID != id ||
				e.AggregateName != name {
				continue
			}

			if !yield(e, nil) {
				return
			}
		}
	}
}

func (mes *MockEventSource) Append(
	ctx context.Context,
	events ...Event,
) error {
	if mes.events == nil {
		mes.events = make([]Event, 0)
	}

	if mes.versions == nil {
		mes.versions = make(map[AggregateID]Version)
	}

	// Snapshot versions
	versions := make(map[AggregateID]Version)

	for _, event := range events {
		if err := event.Valid(); err != nil {
			return err
		}

		if mes.aggregateVersion(versions, event.AggregateID) != event.Version-1 {
			return ErrAggregateVersionMismatch
		}

		versions[event.AggregateID] = event.Version
	}

	// Apply Snapshot
	for id, version := range versions {
		mes.versions[id] = version
	}
	mes.events = append(mes.events, events...)

	return nil
}

func (mes *MockEventSource) aggregateVersion(versions map[AggregateID]Version, id AggregateID) (version Version) {
	a := mes.versions[id]
	b := versions[id]
	if a > b {
		return a
	}
	return b
}
