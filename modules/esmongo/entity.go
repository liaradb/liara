package esmongo

type Entity[I EntityID] interface {
	ID() I
	Version() Version
	Events() []Event
}

type Event interface {
	ID() EventID
	Type() EventType
	EntityID() EntityID
	Version() Version
}

type EntityID interface {
	String() string
}

type Version interface {
	Value() int
}

type EventID interface {
	String() string
}

type EventType interface {
	String() string
}
