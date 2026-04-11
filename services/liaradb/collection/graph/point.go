package graph

import "github.com/liaradb/liaradb/storage/link"

type Point struct {
	id       link.RecordLocator
	edge     link.RecordLocator
	property link.RecordLocator
}

const Pointsize = link.RecordLocatorSize +
	link.RecordLocatorSize +
	link.RecordLocatorSize

func NewPoint(
	id link.RecordLocator,
	edge link.RecordLocator,
	property link.RecordLocator,
) *Point {
	return &Point{
		id:       id,
		edge:     edge,
		property: property,
	}
}

func (p *Point) ID() link.RecordLocator       { return p.id }
func (p *Point) Edge() link.RecordLocator     { return p.edge }
func (p *Point) Property() link.RecordLocator { return p.property }
func (p *Point) Size() int                    { return Pointsize }
