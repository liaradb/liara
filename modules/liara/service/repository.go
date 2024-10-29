package service

import (
	"context"
	"iter"
)

type CommandRepository[I EntityID, E Entity[I]] interface {
	GetByID(context.Context, I) (*E, error)
	Insert(context.Context, *E) error
	Replace(context.Context, *E) error
	Search(context.Context, Query) iter.Seq2[*E, error]
}

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

type Query interface {
	Filter() map[string]any
	Sort() map[string]Sort
	Projection() map[string]Projection
	Offset() int
	Limit() int
}

type Sort int

const (
	SortAsc  Sort = 1
	SortNone Sort = 0
	SortDesc Sort = -1
)

type Projection int

const (
	ProjectionInclude = 1
	ProjectionExclude = 0
)
