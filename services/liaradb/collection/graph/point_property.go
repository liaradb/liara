package graph

import "github.com/liaradb/liaradb/storage/link"

type PointProperty struct {
	id   link.RecordID
	next link.RecordID
}

const PointPropertySize = link.RecordIDSize +
	link.RecordIDSize

func NewPointProperty(
	id link.RecordID,
	next link.RecordID,
) *PointProperty {
	return &PointProperty{
		id:   id,
		next: next,
	}
}

func (pp *PointProperty) ID() link.RecordID   { return pp.id }
func (pp *PointProperty) Next() link.RecordID { return pp.next }
func (pp *PointProperty) Size() int           { return PointPropertySize }
