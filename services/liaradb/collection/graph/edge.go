package graph

import "github.com/liaradb/liaradb/storage/link"

type Edge struct {
	id        link.RecordLocator
	start     link.RecordLocator
	end       link.RecordLocator
	direction Direction
	property  link.RecordLocator
}

const EdgeSize = link.RecordLocatorSize +
	link.RecordLocatorSize +
	link.RecordLocatorSize +
	DirectionSize +
	link.RecordLocatorSize

func NewEdge(
	id link.RecordLocator,
	start link.RecordLocator,
	end link.RecordLocator,
	direction Direction,
	property link.RecordLocator,
) *Edge {
	return &Edge{
		id:        id,
		start:     start,
		end:       end,
		direction: direction,
		property:  property,
	}
}

func (e *Edge) ID() link.RecordLocator       { return e.id }
func (e *Edge) Start() link.RecordLocator    { return e.start }
func (e *Edge) End() link.RecordLocator      { return e.end }
func (e *Edge) Direction() Direction         { return e.direction }
func (e *Edge) Property() link.RecordLocator { return e.property }
func (e *Edge) Size() int                    { return EdgeSize }
