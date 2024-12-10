package liara

import (
	"context"
	"iter"
	"time"
)

type MockEventSource struct {
	requests map[RequestID]mockRequest
	events   []Event
	versions map[AggregateID]Version
}

type mockRequest struct {
	ID   RequestID
	Time time.Time
}

var _ EventRepository = &MockEventSource{}

func (mes *MockEventSource) Get(
	ctx context.Context,
	tenantID TenantID,
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
	tenantID TenantID,
	globalVersion GlobalVersion,
	partitionID []PartitionID,
	limit Limit,
) iter.Seq2[Event, error] {
	panic("unimplemented")
}

func (mes *MockEventSource) GetByOutbox(
	ctx context.Context,
	tenantID TenantID,
	outboxID OutboxID,
	limit Limit,
) iter.Seq2[Event, error] {
	panic("unimplemented")
}

func (mes *MockEventSource) GetByAggregateIDAndName(
	ctx context.Context,
	tenantID TenantID,
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
	tenantID TenantID,
	options AppendOptions,
	events ...AppendEvent,
) error {
	if mes.requests == nil {
		mes.requests = make(map[RequestID]mockRequest)
	}

	if mes.events == nil {
		mes.events = make([]Event, 0)
	}

	if mes.versions == nil {
		mes.versions = make(map[AggregateID]Version)
	}

	if !mes.isNewRequest(options.RequestID) {
		// TODO: Should this be nil?
		return nil
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

	globalVersion := len(mes.events)

	// Apply Snapshot
	for id, version := range versions {
		mes.versions[id] = version
		mes.addRequest(options.RequestID)
	}
	data := make([]Event, 0, len(events))
	for _, event := range events {
		globalVersion++
		data = append(data, mes.toEvent(GlobalVersion(globalVersion), options, event))
	}
	mes.events = append(mes.events, data...)

	return nil
}

func (mes *MockEventSource) isNewRequest(requestID RequestID) bool {
	if requestID == "" {
		return true
	}

	_, ok := mes.requests[requestID]
	return !ok
}

func (mes *MockEventSource) addRequest(requestID RequestID) {
	if requestID == "" {
		return
	}

	if mes.requests == nil {
		mes.requests = make(map[RequestID]mockRequest)
	}

	mes.requests[requestID] = mockRequest{
		ID:   requestID,
		Time: time.Now(),
	}
}

func (mes *MockEventSource) aggregateVersion(versions map[AggregateID]Version, id AggregateID) (version Version) {
	a := mes.versions[id]
	b := versions[id]
	if a > b {
		return a
	}
	return b
}

func (mes *MockEventSource) toEvent(globalVersion GlobalVersion, o AppendOptions, ae AppendEvent) Event {
	return Event{
		GlobalVersion: globalVersion,
		ID:            ae.ID,
		AggregateID:   ae.AggregateID,
		Version:       ae.Version,
		PartitionID:   ae.PartitionID,
		AggregateName: ae.AggregateName,
		Name:          ae.Name,
		Schema:        ae.Schema,
		Metadata: EventMetadata{
			UserID:        o.UserID,
			CorrelationID: o.CorrelationID,
			Time:          o.Time},
		Data: ae.Data,
	}
}
