package eventsource

type (
	Snapshot[T AggregateRoot[U], U ~string] struct {
		AggregateID U
		EventID     EventID
		Version     Version
		Schema      Schema
		Data        T
	}
)
