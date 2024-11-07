package esmongo

type Entity[I EntityID] interface {
	ID() I
	Version() Version
	Events() []Event
}

type EntityID interface {
	String() string
}

type Version interface {
	Value() int
}
