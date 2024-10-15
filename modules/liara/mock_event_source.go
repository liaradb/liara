package liara

import (
	"context"
	"iter"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
)

type MockEventSource struct {
	events   []entity.Event
	versions map[value.AggregateID]value.Version
}

var _ EventRepository = &MockEventSource{}

func (mes *MockEventSource) Get(
	ctx context.Context,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
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
	globalVersion value.GlobalVersion,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	panic("unimplemented")
}

func (mes *MockEventSource) GetByAggregateIDAndName(
	ctx context.Context,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
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
	events ...AppendEvent,
) error {
	if mes.events == nil {
		mes.events = make([]entity.Event, 0)
	}

	if mes.versions == nil {
		mes.versions = make(map[value.AggregateID]value.Version)
	}

	// Snapshot versions
	versions := make(map[value.AggregateID]value.Version)

	for _, event := range events {
		if err := event.Valid(); err != nil {
			return err
		}

		if mes.aggregateVersion(versions, event.AggregateID) != event.Version-1 {
			return value.ErrAggregateVersionMismatch
		}

		versions[event.AggregateID] = event.Version
	}

	globalVersion := len(mes.events)

	// Apply Snapshot
	for id, version := range versions {
		mes.versions[id] = version
	}
	data := make([]entity.Event, 0, len(events))
	for _, event := range events {
		globalVersion++
		data = append(data, mes.toEvent(value.GlobalVersion(globalVersion), event))
	}
	mes.events = append(mes.events, data...)

	return nil
}

func (mes *MockEventSource) aggregateVersion(versions map[value.AggregateID]value.Version, id value.AggregateID) (version value.Version) {
	a := mes.versions[id]
	b := versions[id]
	if a > b {
		return a
	}
	return b
}

func (mes *MockEventSource) toEvent(globalVersion value.GlobalVersion, ae AppendEvent) Event {
	return Event{
		GlobalVersion: globalVersion,
		AggregateName: ae.AggregateName,
		ID:            ae.ID,
		AggregateID:   ae.AggregateID,
		Version:       ae.Version,
		Name:          ae.Name,
		CorrelationID: ae.CorrelationID,
		IdempotenceID: ae.IdempotenceID,
		PartitionID:   ae.PartitionID,
		UserID:        ae.UserID,
		Time:          ae.Time,
		Schema:        ae.Schema,
		Data:          ae.Data,
	}
}
