package esmongo

type Entity[I EntityID] interface {
	ID() I
	Version() int
}

type EntityID interface {
	String() string
}
