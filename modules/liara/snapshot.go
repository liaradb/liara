package liara

import "github.com/cardboardrobots/eventsource/value"

type (
	Snapshot[T AggregateRoot[U], U ~string] struct {
		AggregateID U
		EventID     value.EventID
		Version     value.Version
		Schema      value.Schema
		Data        T
	}
)
