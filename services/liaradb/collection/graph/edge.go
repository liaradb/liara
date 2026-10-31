package graph

import "github.com/liaradb/liaradb/storage/link"

type Edge struct {
	id        link.RecordID
	start     link.RecordID
	end       link.RecordID
	direction Direction
	property  link.RecordID
}

const EdgeSize = link.RecordIDSize +
	link.RecordIDSize +
	link.RecordIDSize +
	DirectionSize +
	link.RecordIDSize

func NewEdge(
	id link.RecordID,
	start link.RecordID,
	end link.RecordID,
	direction Direction,
	property link.RecordID,
) *Edge {
	return &Edge{
		id:        id,
		start:     start,
		end:       end,
		direction: direction,
		property:  property,
	}
}

func (e *Edge) ID() link.RecordID       { return e.id }
func (e *Edge) Start() link.RecordID    { return e.start }
func (e *Edge) End() link.RecordID      { return e.end }
func (e *Edge) Direction() Direction    { return e.direction }
func (e *Edge) Property() link.RecordID { return e.property }
func (e *Edge) Size() int               { return EdgeSize }
