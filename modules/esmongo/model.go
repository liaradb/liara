package esmongo

import "encoding/json"

type model[T any] struct {
	ID      string        `bson:"_id"`
	Version int           `bson:"version"`
	Events  []*modelEvent `bson:"events"`
	Value   T             `bson:"inline"`
}

func newModel[T any](
	id string,
	version int,
	t T,
	events []Event,
) *model[T] {
	evs, _ := newModelEvents(events)
	return &model[T]{
		ID:      id,
		Version: version,
		Events:  evs,
		Value:   t,
	}
}

func (m *model[T]) increment() *model[T] {
	m.Version++
	return m
}

type modelEvent struct {
	Type string `bson:"type"`
	Data []byte `bson:"data"`
}

func newModelEvent(
	eventType string,
	e any,
) (*modelEvent, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &modelEvent{
		Type: eventType,
		Data: data,
	}, nil
}

func newModelEvents(
	events []Event,
) ([]*modelEvent, error) {
	result := make([]*modelEvent, 0, len(events))

	for _, e := range events {
		r, err := newModelEvent(
			e.Type(),
			e)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}
