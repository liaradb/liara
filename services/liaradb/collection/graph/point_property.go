package graph

import "github.com/liaradb/liaradb/storage/link"

type PointProperty struct {
	id   link.RecordLocator
	next link.RecordLocator
}

const PointPropertySize = link.RecordLocatorSize +
	link.RecordLocatorSize

func NewPointProperty(
	id link.RecordLocator,
	next link.RecordLocator,
) *PointProperty {
	return &PointProperty{
		id:   id,
		next: next,
	}
}

func (pp *PointProperty) ID() link.RecordLocator   { return pp.id }
func (pp *PointProperty) Next() link.RecordLocator { return pp.next }
func (pp *PointProperty) Size() int                { return PointPropertySize }
