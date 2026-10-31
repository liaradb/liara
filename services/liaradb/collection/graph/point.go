package graph

import "github.com/liaradb/liaradb/storage/link"

type Point struct {
	id       link.RecordID
	edge     link.RecordID
	property link.RecordID
}

const Pointsize = link.RecordIDSize +
	link.RecordIDSize +
	link.RecordIDSize

func NewPoint(
	id link.RecordID,
	edge link.RecordID,
	property link.RecordID,
) *Point {
	return &Point{
		id:       id,
		edge:     edge,
		property: property,
	}
}

func (p *Point) ID() link.RecordID       { return p.id }
func (p *Point) Edge() link.RecordID     { return p.edge }
func (p *Point) Property() link.RecordID { return p.property }
func (p *Point) Size() int               { return Pointsize }
