package esmongo

type Event interface {
	ID() EventID
	Type() EventType
	EntityID() EntityID
	Version() Version
}

type EventID interface {
	String() string
}

type EventType interface {
	String() string
}
