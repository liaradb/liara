package esmongo

type model[T any] struct {
	ID      string        `bson:"_id"`
	Version int           `bson:"_version"`
	Schema  string        `bson:"_schema"`
	Events  []*modelEvent `bson:"_events"`
	Value   T             `bson:"inline"`
}

func newModel[T any](
	id string,
	version int,
	schema string,
	events []*modelEvent,
	value T,
) *model[T] {
	return &model[T]{
		ID:      id,
		Version: version,
		Schema:  schema,
		Events:  events,
		Value:   value,
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
