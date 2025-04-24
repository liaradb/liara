package esmongo

type model[T any] struct {
	ID      string        `bson:"_id"`
	Version int           `bson:"_version"`
	Events  []*modelEvent `bson:"_events"`
	Value   T             `bson:"inline"`
}

func newModel[T any](
	id string,
	version int,
	t T,
	events []*modelEvent,
) *model[T] {
	return &model[T]{
		ID:      id,
		Version: version,
		Events:  events,
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
	data []byte,
) *modelEvent {
	return &modelEvent{
		Type: eventType,
		Data: data,
	}
}
